/*
Copyright 2018 The Rook Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package mgr for the Ceph manager.
package mgr

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/koor-tech/koor/pkg/daemon/ceph/client"
	"github.com/koor-tech/koor/pkg/util"
	"github.com/koor-tech/koor/pkg/util/exec"
	"github.com/pkg/errors"
)

const (
	dashboardUserReadOnlyRole = "read-only"
)

type DashboardSSOInfo struct {
	OneLoginSettings OneLoginSettings `json:"onelogin_settings"`
}

type OneLoginSettings struct {
	// TODO Check which fields are important to check against the CRD input
}

func (c *Cluster) configureSSO() (bool, error) {
	if !c.spec.Dashboard.SSO.Enabled {
		// Check if SSO is still enabled
		args := []string{"dashboard", "sso", "status"}
		for i := 0; i < 5; i++ {
			out, err := client.NewCephCommand(c.context, c.clusterInfo, args).RunWithTimeout(exec.CephCommandsTimeout)
			if err == context.DeadlineExceeded {
				logger.Warning("SSO disable timed out. trying again")
				continue
			}
			if strings.Contains(string(out), "disabled") {
				return false, nil
			}

			logger.Infof("Disabling dashboard SSO")
			// Make sure SSO is disabled
			args = []string{"dashboard", "sso", "disable"}
			for i := 0; i < 5; i++ {
				_, err := client.NewCephCommand(c.context, c.clusterInfo, args).RunWithTimeout(exec.CephCommandsTimeout)
				if err == context.DeadlineExceeded {
					logger.Warning("SSO disable timed out. trying again")
					continue
				}
				return true, nil
			}
		}

		return true, errors.New("failed to disable SSO")
	}

	logger.Infof("Enabling dashboard SSO")

	// create and build sso setup args command
	dashboardURL := c.spec.Dashboard.SSO.BaseURL
	idpMetadataURL := c.spec.Dashboard.SSO.IDPMetadataURL
	idpUsernameAttribute := c.spec.Dashboard.SSO.IDPAttributes.Username
	idpEntityID := c.spec.Dashboard.SSO.EntityID

	// TODO check the `ceph dashboard sso show saml2` json output if we need to re-setup the SAML2 config in the dashboard
	ssout, err := client.NewCephCommand(c.context, c.clusterInfo, []string{"dashboard", "sso", "show", "saml2"}).RunWithTimeout(exec.CephCommandsTimeout)
	if err != nil {
		return false, err
	}
	_ = ssout

	args := []string{"dashboard", "sso", "setup", "saml2", dashboardURL, idpMetadataURL}
	if idpUsernameAttribute != "" || idpEntityID != "" {
		args = append(args, idpUsernameAttribute, idpEntityID)
	}

	// retry a few times in the case that the mgr module is not ready to accept commands
	for i := 0; i < 5; i++ {
		_, err := client.NewCephCommand(c.context, c.clusterInfo, args).RunWithTimeout(exec.CephCommandsTimeout)
		if err == nil {
			break
		}
		if err == context.DeadlineExceeded {
			logger.Warning("sso setup timed out. trying again")
			continue
		}

		exitCode, parsed := c.exitCode(err)
		if parsed {
			if exitCode == invalidArgErrorCode {
				logger.Info("dashboard module is not ready yet. trying again")
				time.Sleep(dashboardInitWaitTime)
				continue
			}
		}
		return false, errors.Wrap(err, "failed to setup sso on mgr dashboard")
	}

	return c.createUsers()
}

type dashboardUserInfo struct {
	Roles []string `json:"roles"`
}

func (c *Cluster) createUsers() (bool, error) {
	// Generate a password
	password, err := GeneratePassword(passwordLength)
	if err != nil {
		return false, errors.Wrap(err, "failed to generate password")
	}

	file, err := util.CreateTempFile(password)
	if err != nil {
		return false, errors.Wrap(err, "failed to create a temporary dashboard password file")
	}
	defer func() {
		if err := os.Remove(file.Name()); err != nil {
			logger.Errorf("failed to clean up dashboard password file %q. %v", file.Name(), err)
		}
	}()

	args := []string{"dashboard", "ac-user-show"}
	userOutput, err := client.NewCephCommand(c.context, c.clusterInfo, args).RunWithTimeout(exec.CephCommandsTimeout)
	if err != nil {
		return false, errors.Wrap(err, "failed to get users")
	}
	var usersParsed []string
	if err := json.Unmarshal(userOutput, &usersParsed); err != nil {
		return false, errors.Wrap(err, "failed to parse dashboard user list")
	}
	users := map[string]interface{}{}
	for _, user := range usersParsed {
		users[user] = nil
	}

	var changed bool
	for _, user := range c.spec.Dashboard.SSO.Users {
		if len(user.Roles) == 0 {
			user.Roles = append(user.Roles, dashboardUserReadOnlyRole)
		}

		// If the user already exists we make sure the roles are set accordingly
		if _, ok := users[user.Username]; !ok {
			args := []string{"dashboard", "ac-user-create", user.Username, "-i", file.Name(), user.Roles[0]}
			_, err = client.ExecuteCephCommandWithRetry(func() (string, []byte, error) {
				output, err := client.NewCephCommand(c.context, c.clusterInfo, args).RunWithTimeout(exec.CephCommandsTimeout)
				return "create dashboard user", output, err
			}, c.exitCode, 5, invalidArgErrorCode, dashboardInitWaitTime)
			if err != nil {
				return false, errors.Wrap(err, "failed to create user")
			}
		}

		// If the user already exists, update the role
		if changed, err = c.ensureUserRoles(user.Username, user.Roles); err != nil {
			return false, errors.Wrap(err, "")
		}
	}

	logger.Info("successfully created dashboard sso users")
	return changed, nil
}

func (c *Cluster) getUserRoles(username string) ([]string, error) {
	argsRoles := []string{"dashboard", "ac-user-show", username}
	roleOutput, err := client.NewCephCommand(c.context, c.clusterInfo, argsRoles).RunWithTimeout(exec.CephCommandsTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user and role info")
	}
	var userInfo dashboardUserInfo
	if err := json.Unmarshal(roleOutput, &userInfo); err != nil {
		return nil, errors.Wrap(err, "failed to parse dashboard user info")
	}

	return userInfo.Roles, nil
}

func (c *Cluster) ensureUserRoles(username string, roles []string) (bool, error) {
	currentRoles, err := c.getUserRoles(username)
	if err != nil {
		return false, errors.Wrap(err, "failed to ")
	}

	// Make the current roles into a map so we just check if the key is set
	shouldHaveRoles := map[string]interface{}{}
	for _, role := range roles {
		shouldHaveRoles[role] = nil
	}
	for _, role := range currentRoles {
		delete(shouldHaveRoles, role)
	}
	// Convert back the map to a string slice
	rolesToAdd := make([]string, len(shouldHaveRoles))
	i := 0
	for role := range shouldHaveRoles {
		rolesToAdd[i] = role
		i++
	}

	// Check if we need to set user roles
	if len(rolesToAdd) == 0 {
		return false, nil
	}

	if err := c.setUserRoles(username, roles); err != nil {
		return false, errors.Wrap(err, "failed to remove user roles")
	}

	return true, nil
}

func (c *Cluster) setUserRoles(username string, roles []string) error {
	args := []string{"dashboard", "ac-user-set-roles", username}
	args = append(args, roles...)
	_, err := client.ExecuteCephCommandWithRetry(func() (string, []byte, error) {
		output, err := client.NewCephCommand(c.context, c.clusterInfo, args).RunWithTimeout(exec.CephCommandsTimeout)
		return "set user role", output, err
	}, c.exitCode, 5, invalidArgErrorCode, dashboardInitWaitTime)
	if err != nil {
		return errors.Wrap(err, "failed to set user roles")
	}

	return nil
}

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
	"encoding/json"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/rook/rook/pkg/daemon/ceph/client"
	"github.com/rook/rook/pkg/util"
	"github.com/rook/rook/pkg/util/exec"
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
		output, err := client.ExecuteCephCommandWithRetry(func() (string, []byte, error) {
			output, err := client.NewCephCommand(c.context, c.clusterInfo, args).RunWithTimeout(exec.CephCommandsTimeout)
			return "checking dashboard sso status", output, err
		}, c.exitCode, 5, invalidArgErrorCode, dashboardInitWaitTime)
		if err != nil {
			return false, errors.Wrap(err, "failed to check dashboard sso status")
		}
		// SSO is already disabled, no need to disable it (again)
		if strings.Contains(string(output), "disabled") {
			logger.Infof("Dashboard SSO already disabled")
			return false, nil
		}

		logger.Infof("Disabling dashboard SSO")
		// Make sure SSO is disabled
		args = []string{"dashboard", "sso", "disable"}
		_, err = client.ExecuteCephCommandWithRetry(func() (string, []byte, error) {
			output, err := client.NewCephCommand(c.context, c.clusterInfo, args).RunWithTimeout(exec.CephCommandsTimeout)
			return "disable dashboard sso", output, err
		}, c.exitCode, 5, invalidArgErrorCode, dashboardInitWaitTime)
		return true, err
	}

	logger.Infof("Enabling dashboard SSO")

	// create and build sso setup args command
	dashboardURL := c.spec.Dashboard.SSO.BaseURL
	idpMetadataURL := c.spec.Dashboard.SSO.IDPMetadataURL
	idpAttributes := c.spec.Dashboard.SSO.IDPAttributes
	idpEntityID := c.spec.Dashboard.SSO.EntityID

	// TODO check the `ceph dashboard sso show saml2` json output if we need to re-setup the SAML2 config in the dashboard
	_, err := client.NewCephCommand(c.context, c.clusterInfo, []string{"dashboard", "sso", "show", "saml2"}).RunWithTimeout(exec.CephCommandsTimeout)
	if err != nil {
		return false, err
	}

	args := []string{"dashboard", "sso", "setup", "saml2", dashboardURL, idpMetadataURL}
	if idpAttributes.Username != "" || idpEntityID != "" {
		args = append(args, idpAttributes.Username, idpEntityID)
	}

	// retry a few times in the case that the mgr module is not ready to accept commands
	_, err = client.ExecuteCephCommandWithRetry(func() (string, []byte, error) {
		output, err := client.NewCephCommand(c.context, c.clusterInfo, args).RunWithTimeout(exec.CephCommandsTimeout)
		return "setup dashboard sso", output, err
	}, c.exitCode, 5, invalidArgErrorCode, dashboardInitWaitTime)
	if err != nil {
		return false, err
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

		// If the user doesn't exist, we create it and make sure the roles are set accordingly
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
	logger.Info("SSO Setup Successful")
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

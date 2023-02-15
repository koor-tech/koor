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
	"os"
	"time"

	"github.com/koor-tech/koor/pkg/daemon/ceph/client"
	"github.com/koor-tech/koor/pkg/util"
	"github.com/koor-tech/koor/pkg/util/exec"
	"github.com/pkg/errors"
)

func (c *Cluster) setupSSO() (bool, error) {
	logger.Infof("Starting SSO Setup: enabled=%+v", c.spec.Dashboard.SSO.Enabled)
	if !c.spec.Dashboard.SSO.Enabled {
		// Make sure SSO is disabled
		args := []string{"dashboard", "sso", "disable"}
		for i := 0; i < 5; i++ {
			_, err := client.NewCephCommand(c.context, c.clusterInfo, args).RunWithTimeout(exec.CephCommandsTimeout)
			if err == context.DeadlineExceeded {
				logger.Warning("SSO disable timed out. trying again")
			}
		}
		return false, nil
	}

	// create and build sso setup args command
	dashboardURL := c.spec.Dashboard.SSO.BaseURL
	idpMetadataURL := c.spec.Dashboard.SSO.IDPMetadataURL
	idpUsernameAttribute := c.spec.Dashboard.SSO.IDPAttributes.Username
	idpEntityID := c.spec.Dashboard.SSO.EntityID

	// TODO this should check if the configuration has changed because otherwise if
	// SSO is enabled once, it must be disabled first and then re-enabled with the new settings
	// run `ceph dashboard sso show saml2` and search for the input args

	ssout, err := client.NewCephCommand(c.context, c.clusterInfo, []string{"dashboard", "sso", "show", "saml2"}).RunWithTimeout(exec.CephCommandsTimeout)
	if err == nil {
		return false, err
	}
	// TODO Check the output for the settings if they are all still set the same (simple grep'ing should be enough)

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

	return true, nil
}

// TODO add function to create the users and set them the accordingly role
// TODO if an user already exists, the user needs to be checked to only have the roles as specified in the list

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

	users := c.spec.Dashboard.SSO.Users
	for i := 0; i < len(users); i++ {
		username := users[i].Username
		role := users[i].Role
		args := []string{"dashboard", "ac-user-create", username, "-i", file.Name(), role}
		defer func() {
			if err := os.Remove(file.Name()); err != nil {
				logger.Errorf("failed to clean up dashboard password file %q. %v", file.Name(), err)
			}
		}()
		_, err = client.ExecuteCephCommandWithRetry(func() (string, []byte, error) {
			output, err := client.NewCephCommand(c.context, c.clusterInfo, args).RunWithTimeout(exec.CephCommandsTimeout)
			return "create dashboard user", output, err
		}, c.exitCode, 5, invalidArgErrorCode, dashboardInitWaitTime)
		if err != nil {
			return false, errors.Wrap(err, "failed to create user")
		}
		// If the user already exists, update the role
		if err.Error() == "User"+"'"+username+"' already exists" {
			args := []string{"dashboard", "ac-user-set-roles", username, role}
			_, err = client.ExecuteCephCommandWithRetry(func() (string, []byte, error) {
				output, err := client.NewCephCommand(c.context, c.clusterInfo, args).RunWithTimeout(exec.CephCommandsTimeout)
				return "set user role", output, err
			}, c.exitCode, 5, invalidArgErrorCode, dashboardInitWaitTime)
			if err != nil {
				return false, errors.Wrap(err, "failed to update role")
			}
		}

		logger.Info("successfully created dashboard user")
		return true, nil
	}
}

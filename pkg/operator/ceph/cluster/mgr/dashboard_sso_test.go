/*
Copyright (c) 2023 Koor Technolgies, Inc.
*/
package mgr

import (
	"context"
	"strings"
	"testing"
	"time"

	cephv1 "github.com/rook/rook/pkg/apis/ceph.rook.io/v1"
	"github.com/rook/rook/pkg/clusterd"
	cephclient "github.com/rook/rook/pkg/daemon/ceph/client"
	"github.com/rook/rook/pkg/operator/test"
	exectest "github.com/rook/rook/pkg/util/exec/test"
	"github.com/stretchr/testify/assert"
)

func TestEnableSSO(t *testing.T) {
	ctx := context.TODO()
	statuses := 0
	shows := 0
	setups := 0
	disables := 0
	userCreates := 0
	userSetRoles := 0
	exitCodeResponse := 0
	clientset := test.New(t, 3)
	mockFN := func(command string, args ...string) (string, error) {
		logger.Infof("command: %s %v", command, args)
		exitCodeResponse = 0
		if args[0] == "dashboard" {
			if args[1] == "sso" {
				if args[2] == "status" {
					if statuses > 0 {
						return "disabled", nil
					}
					statuses++
					return "enabled", nil
				} else if args[2] == "show" {
					shows++
					return "{\"onelogin_settings\":{}}", nil
				} else if args[2] == "setup" {
					setups++
				} else if args[2] == "disable" {
					disables++
					return "SSO is \"disabled\".", nil
				}
			} else if strings.HasPrefix(args[1], "ac-user-") {
				if args[1] == "ac-user-show" {
					if args[2] == "example123" {
						if userSetRoles > 0 {
							return "{\"roles\":[\"block-manager\", \"rgw-manager\"]}", nil
						} else {
							return "{\"roles\":[\"block-manager\"]}", nil
						}
					}
					if userCreates > 0 {
						return "[\"admin\", \"example123\"]", nil
					} else {
						return "[\"admin\"]", nil
					}
				} else if args[1] == "ac-user-create" {
					userCreates++
				} else if args[1] == "ac-user-set-roles" {
					userSetRoles++
				}
			}
		}
		return "", nil
	}
	executor := &exectest.MockExecutor{
		MockExecuteCommandWithOutput: mockFN,
		MockExecuteCommandWithTimeout: func(timeout time.Duration, command string, arg ...string) (string, error) {
			return mockFN(command, arg...)
		},
	}

	ownerInfo := cephclient.NewMinimumOwnerInfoWithOwnerRef()
	clusterInfo := &cephclient.ClusterInfo{
		Namespace: "myns",
		OwnerInfo: ownerInfo,
		Context:   ctx,
	}
	c := &Cluster{clusterInfo: clusterInfo, context: &clusterd.Context{Clientset: clientset, Executor: executor},
		spec: cephv1.ClusterSpec{
			Dashboard: cephv1.DashboardSpec{SSO: cephv1.SSOSpec{
				Enabled: true,
				Users: []cephv1.UserRef{
					{
						Username: "example123",
						Roles: []string{
							"block-manager",
							"rgw-manager",
						},
					},
				},
			}},
		},
	}
	c.exitCode = func(err error) (int, bool) {
		if exitCodeResponse != 0 {
			return exitCodeResponse, true
		}
		return exitCodeResponse, false
	}

	dashboardInitWaitTime = 0
	// enable SSO
	changed, err := c.configureSSO()
	assert.True(t, changed)
	assert.NoError(t, err)
	assert.Equal(t, 1, setups)
	assert.Equal(t, 0, disables)
	assert.Equal(t, 1, userCreates)
	assert.Equal(t, 1, userSetRoles)

	// enabled SSO should be a no-op on second runs
	changed, err = c.configureSSO()
	assert.False(t, changed)
	assert.Nil(t, err)
	assert.Equal(t, 2, setups)
	assert.Equal(t, 0, disables)
	assert.Equal(t, 1, userCreates)
	assert.Equal(t, 1, userSetRoles)

	// disable SSO
	c.spec.Dashboard.SSO.Enabled = false
	changed, err = c.configureSSO()
	assert.True(t, changed)
	assert.Nil(t, err)
	assert.Equal(t, 2, setups)
	assert.Equal(t, 1, disables)
	assert.Equal(t, 1, userCreates)
	assert.Equal(t, 1, userSetRoles)

	// disabled SSO should be a no-op on second runs
	changed, err = c.configureSSO()
	assert.False(t, changed)
	assert.Nil(t, err)
	assert.Equal(t, 2, setups)
	assert.Equal(t, 1, disables)
	assert.Equal(t, 1, userCreates)
	assert.Equal(t, 1, userSetRoles)
}

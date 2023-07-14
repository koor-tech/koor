/*
Copyright (c) 2023 Koor Technolgies, Inc.
*/

package osd

import (
	"context"
	"testing"
	"time"

	cephv1 "github.com/rook/rook/pkg/apis/ceph.rook.io/v1"
	"github.com/rook/rook/pkg/clusterd"
	cephclient "github.com/rook/rook/pkg/daemon/ceph/client"
	"github.com/rook/rook/pkg/operator/ceph/config"
	cephver "github.com/rook/rook/pkg/operator/ceph/version"
	exectest "github.com/rook/rook/pkg/util/exec/test"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestScrubbingScheduleSetup(t *testing.T) {
	executor := &exectest.MockExecutor{}
	clientset := fake.NewSimpleClientset()
	clusterInfo := &cephclient.ClusterInfo{
		Namespace:   "ns",
		CephVersion: cephver.Quincy,
		Context:     context.TODO(),
		OwnerInfo:   cephclient.NewMinimumOwnerInfoWithOwnerRef(),
	}
	ctx := &clusterd.Context{
		Clientset: clientset,
		Executor:  executor}
	spec := cephv1.ClusterSpec{}
	c := New(ctx, clusterInfo, spec, "myversion")
	s := config.GetMonStore(ctx, clusterInfo)

	execCount := 0
	countExecs := true
	mockStore := map[string]map[string]string{}
	executor.MockExecuteCommandWithTimeout =
		func(timeout time.Duration, command string, args ...string) (string, error) {
			if countExecs {
				execCount++
			}

			if args[0] == "config" {
				if args[1] == "get" {
					whose, ok := mockStore[args[2]]
					if !ok {
						return "", nil
					}

					val, ok := whose[args[3]]
					if !ok {
						return "", nil
					}

					return val, nil
				} else if args[1] == "set" {
					whose, ok := mockStore[args[2]]
					if !ok {
						whose = map[string]string{}
					}

					whose[args[3]] = args[4]
					mockStore[args[2]] = whose

					return "", nil
				}
			}

			return "", nil
		}

	// No scrubbing schedule given, so no commands are executed
	err := c.ConfigureOSDScrubbing()
	assert.Nil(t, err)
	assert.Equal(t, 0, execCount)

	err = c.ConfigureOSDScrubbing()
	assert.Nil(t, err)
	assert.Equal(t, 0, execCount)

	countExecs = false
	emptyExpect := map[string]string{
		"osd_max_scrubs":           "",
		"osd_scrub_begin_hour":     "",
		"osd_scrub_end_hour":       "",
		"osd_scrub_begin_week_day": "",
		"osd_scrub_end_week_day":   "",
		"osd_scrub_min_interval":   "",
		"osd_scrub_max_interval":   "",
		"osd_deep_scrub_interval":  "",
		"osd_scrub_sleep":          "",
	}
	validateConfigStoreContent(t, s, "osd", emptyExpect)
	countExecs = true

	// Scrubbing schedule is enabled, it should be configured now
	// Every config option causes 1 get to the config store
	// Every set causes 1 set to the config store
	schedule := cephv1.Scrubbing{
		ApplySchedule:     true,
		MaxScrubOps:       2,
		BeginHour:         21,
		EndHour:           5,
		BeginWeekDay:      1,
		EndWeekDay:        5,
		MinScrubInterval:  &metav1.Duration{Duration: defaultScrubMinInterval},
		MaxScrubInterval:  &metav1.Duration{Duration: defaultScrubMaxInterval},
		DeepScrubInterval: &metav1.Duration{Duration: defaultDeepScrubInterval},
		ScrubSleepSeconds: &metav1.Duration{Duration: defaultScrubSleepSeconds},
	}
	c.spec.Storage.Scrubbing = schedule
	err = c.ConfigureOSDScrubbing()
	assert.Nil(t, err)
	assert.Equal(t, 18, execCount)

	err = c.ConfigureOSDScrubbing()
	assert.Nil(t, err)
	assert.Equal(t, 27, execCount)

	countExecs = false
	expected := map[string]string{
		"osd_max_scrubs":           "2",
		"osd_scrub_begin_hour":     "21",
		"osd_scrub_end_hour":       "5",
		"osd_scrub_begin_week_day": "1",
		"osd_scrub_end_week_day":   "5",
		"osd_scrub_min_interval":   "86400.000000",
		"osd_scrub_max_interval":   "604800.000000",
		"osd_deep_scrub_interval":  "604800.000000",
		"osd_scrub_sleep":          "0.000000",
	}
	validateConfigStoreContent(t, s, "osd", expected)
	countExecs = true

	// Disable scrubbing and not changes should occur
	c.spec.Storage.Scrubbing.ApplySchedule = false

	err = c.ConfigureOSDScrubbing()
	assert.Nil(t, err)
	assert.Equal(t, 27, execCount)

	err = c.ConfigureOSDScrubbing()
	assert.Nil(t, err)
	assert.Equal(t, 27, execCount)

	// Enable and change a value of the scrubbing schedule
	c.spec.Storage.Scrubbing.ApplySchedule = true
	c.spec.Storage.Scrubbing.BeginHour = 1
	expected["osd_scrub_begin_hour"] = "1"
	c.spec.Storage.Scrubbing.ScrubSleepSeconds = &metav1.Duration{Duration: 10*time.Second + 500*time.Millisecond}
	expected["osd_scrub_sleep"] = "10.500000"

	err = c.ConfigureOSDScrubbing()
	assert.Nil(t, err)
	assert.Equal(t, 38, execCount)

	err = c.ConfigureOSDScrubbing()
	assert.Nil(t, err)
	assert.Equal(t, 47, execCount)

	countExecs = false
	validateConfigStoreContent(t, s, "osd", expected)
}

func validateConfigStoreContent(t *testing.T, s *config.MonStore, who string, expected map[string]string) {
	for key, expect := range expected {
		val, err := s.Get(who, key)
		assert.Nil(t, err)
		assert.Equalf(t, expect, val, "key: %s", key)
	}
}

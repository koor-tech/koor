/*
Copyright (c) 2023 Koor Technolgies, Inc.
*/

package osd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rook/rook/pkg/operator/ceph/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaultScrubMinInterval  = 1 * 24 * time.Hour
	defaultScrubMaxInterval  = 7 * 24 * time.Hour
	defaultDeepScrubInterval = 7 * 24 * time.Hour
	defaultScrubSleepSeconds = 0 * time.Second
)

func (c *Cluster) ConfigureOSDScrubbing() error {
	if !c.spec.Storage.Scrubbing.ApplySchedule {
		logger.Info("scrubbing schedule not enabled, skipping setting custom schedule")
		return nil
	}

	logger.Debug("applying scrubbing schedule")

	schedule := c.spec.Storage.Scrubbing
	s := config.GetMonStore(c.context, c.clusterInfo)

	// Max Scrubs
	if schedule.MaxScrubOps < 1 {
		schedule.MaxScrubOps = 3
		logger.Warning("invalid schrubbing max scrub ops set, must be 1 or higher")
	}
	if _, err := s.SetIfChanged("osd", "osd_max_scrubs", strconv.Itoa(int(schedule.MaxScrubOps))); err != nil {
		return err
	}

	// Hour
	if _, err := s.SetIfChanged("osd", "osd_scrub_begin_hour", strconv.Itoa(int(schedule.BeginHour))); err != nil {
		return err
	}
	if _, err := s.SetIfChanged("osd", "osd_scrub_end_hour", strconv.Itoa(int(schedule.EndHour))); err != nil {
		return err
	}

	// Week Days
	if _, err := s.SetIfChanged("osd", "osd_scrub_begin_week_day", strconv.Itoa(int(schedule.BeginWeekDay))); err != nil {
		return err
	}
	if _, err := s.SetIfChanged("osd", "osd_scrub_end_week_day", strconv.Itoa(int(schedule.EndWeekDay))); err != nil {
		return err
	}

	// Scrub Intervals
	if schedule.MinScrubInterval == nil || schedule.MinScrubInterval.Duration <= 0 {
		schedule.MinScrubInterval.Duration = defaultScrubMinInterval
	}
	if _, err := s.SetIfChanged("osd", "osd_scrub_min_interval", fmt.Sprintf("%f", schedule.MinScrubInterval.Seconds())); err != nil {
		return err
	}
	if schedule.MaxScrubInterval == nil || schedule.MaxScrubInterval.Duration <= 0 {
		schedule.MaxScrubInterval.Duration = defaultScrubMaxInterval
	}
	if _, err := s.SetIfChanged("osd", "osd_scrub_max_interval", fmt.Sprintf("%f", schedule.MaxScrubInterval.Seconds())); err != nil {
		return err
	}

	if schedule.DeepScrubInterval == nil || schedule.DeepScrubInterval.Duration <= 0 {
		schedule.DeepScrubInterval.Duration = defaultDeepScrubInterval
	}
	if _, err := s.SetIfChanged("osd", "osd_deep_scrub_interval", fmt.Sprintf("%f", schedule.DeepScrubInterval.Seconds())); err != nil {
		return err
	}

	if schedule.ScrubSleepSeconds == nil || schedule.ScrubSleepSeconds.Duration < 0 {
		scrubSleepSeconds := 0 * time.Second
		schedule.ScrubSleepSeconds = &metav1.Duration{
			Duration: scrubSleepSeconds,
		}
	}
	if _, err := s.SetIfChanged("osd", "osd_scrub_sleep", fmt.Sprintf("%f", schedule.ScrubSleepSeconds.Seconds())); err != nil {
		return err
	}

	logger.Info("scrubbing schedule applied if needed")

	return nil
}

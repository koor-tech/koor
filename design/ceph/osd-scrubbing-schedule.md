# OSD Scrubbing Schedule

## Goal

Adding a simple way to configure the Ceph OSD Scrubbing schedule options from the CephCluster CRD.

## Current Solution

Users need to update the `rook-config-override` and add the [Ceph OSD Scrubbing config options](https://docs.ceph.com/en/latest/rados/configuration/osd-config-ref/#scrubbing) as needed.

## Proposed Solution

Adding a new config structure to the CephCluster CRD under `.spec.storage` named `scrubbing:`.

The config structure will contain the important parts for setting a schedule for the OSD scrubbing.

```yaml
spec:
  # [...]
  storage:
    # [...]
    scrubbing:
      applySchedule: true
      maxScrubOps: 3
      beginHour: 8
      endHour: 17
      beginWeekDay: 1
      endWeekDay: 5
      minScrubInterval: 24h
      maxScrubInterval: 168h
      deepScrubInterval: 168h
      scrubSleepSeconds: 100ms
  # [...]
```

This will use the Ceph config store to apply these settings to the cluster.

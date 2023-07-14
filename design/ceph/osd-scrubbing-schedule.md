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
    scrubbingSchedule:
      maxScrubsOps: 3
      beginHour: 8
      beginWeekDay: 1
      endHour: 17
      endWeekDay: 5
      minScrubInterval: 1d
      maxScrubInterval: 7d
      deepScrubInterval: 7d
  # [...]
```

This will use the Ceph config store to apply these settings to the cluster.

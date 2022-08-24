---
title: Requirements
---

## Software

Check out the [prerequisites](../Getting-Started/Prerequisites/prerequisites.md) section to get started.

## Minimum Requirements for Ceph Storage node

### RAM

- 16 GB baseline, plus 5 GB for every OSD on the node

### Network

- For High-performance, high-endurance enterprise NVMe SSDs:
    - 10 Gigabit Ethernet (GbE) per 2 OSDs.
- For HDDs:
    - 10 GbE per 12 OSDs each for client- and cluster-facing networks.

### CPU

- 1 core per HDD OSD minimum

#### For IOPS optimized solutions where there are NVMe SSDs

- 6 cores per NVMe SSD

### Block.db sizing for OSD

The general recommendation is to have block.db size in between 1% to 4% of block size. For RGW workloads, it is recommended that the block.db size isn’t smaller than 4% of block, because RGW heavily uses it to store metadata (omap keys). For example, if the block size is 1TB, then block.db shouldn’t be less than 40GB. For RBD workloads, 1% to 2% of block size is usually enough.

### Block.wal sizing

If there is only a small amount of fast storage available (e.g., less than a gigabyte), we recommend using it as a WAL device. If there is more, provisioning a DB device makes more sense. The BlueStore journal will always be placed on the fastest device available, so using a DB device will provide the same benefit that the WAL device would while also allowing additional metadata to be stored there (if it will fit). This means that if a DB device is specified but an explicit WAL device is not, the WAL will be implicitly colocated with the DB on the faster device.

  *When not using a mix of fast and slow devices, it isn’t required to create separate logical volumes for block.db (or block.wal). BlueStore will automatically colocate these within the space of block.*

### Partition Layout

Either a separate `/var` partition or the OS Root partition should be at least 100 GB for MON nodes.

**OS Disk**:

- At least 100-200GB in total size
- SATA SSD or NVMe SSD.

**Recommended layout**:

- `/` - 20GB
- `/var` - 75-100GB or more
- `/var/lib/containers` - (this assumes CRI-O is used) 75-100GB

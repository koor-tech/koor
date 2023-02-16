---
title: Storage Architecture
---

Ceph is a highly scalable distributed storage solution for **block storage**, **object storage**, and **shared filesystems** with years of production deployments.

## Design

Koor Storage Distribution enables Ceph storage to run on Kubernetes using Kubernetes primitives.
With Ceph running in the Kubernetes cluster, Kubernetes applications can
mount block devices and filesystems managed by Koor Storage Distribution, or can use the S3/Swift API for object storage. The Koor Storage Distribution operator
automates configuration of storage components and monitors the cluster to ensure the storage remains available
and healthy.

The Koor Storage Distribution operator is a simple container that has all that is needed to bootstrap
and monitor the storage cluster. The operator will start and monitor [Ceph monitor pods](../Storage-Configuration/Advanced/ceph-mon-health.md), the Ceph OSD daemons to provide RADOS storage, as well as start and manage other Ceph daemons. The operator manages CRDs for pools, object stores (S3/Swift), and filesystems by initializing the pods and other resources necessary to run the services.

The operator will monitor the storage daemons to ensure the cluster is healthy. Ceph mons will be started or failed over when necessary, and
other adjustments are made as the cluster grows or shrinks.  The operator will also watch for desired state changes
specified in the Ceph custom resources (CRs) and apply the changes.

Koor Storage Distribution automatically configures the Ceph-CSI driver to mount the storage to your pods.

![Koor Storage Distribution Components on Kubernetes](ceph-storage/kubernetes.png)

The `koorinc/ceph` image includes all necessary tools to manage the cluster. Koor Storage Distribution is not in the Ceph data path.
Many of the Ceph concepts like placement groups and crush maps
are hidden so you don't have to worry about them. Instead Koor Storage Distribution creates a simplified user experience for admins that is in terms
of physical resources, pools, volumes, filesystems, and buckets. At the same time, advanced configuration can be applied when needed with the Ceph tools.

Koor Storage Distribution is implemented in golang. Ceph is implemented in C++ where the data path is highly optimized. We believe
this combination offers the best of both worlds.

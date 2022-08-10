---
title: Helm Charts Overview
---

Rook has published the following Helm charts for the Koor Storage Distribution (Ceph):

* [Koor Operator](operator-chart.md): Starts the Ceph Operator, which will watch for Ceph CRs (custom resources)
* [Koor Cluster](ceph-cluster-chart.md): Creates Ceph CRs that the operator will use to configure the cluster

The Helm charts are intended to simplify deployment and upgrades.
Configuring the Rook resources without Helm is also fully supported by creating the
[manifests](https://github.com/koor-tech/koor/tree/master/deploy/examples)
directly.

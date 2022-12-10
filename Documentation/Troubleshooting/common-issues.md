---
title: Common Issues
---

To help troubleshoot your Koor Storage Distribution clusters, here are some tips on what information will help solve the issues you might be seeing.
If after trying the suggestions found on this page and the problem is not resolved, the Koor Storage Distribution team is very happy to help you troubleshoot the issues through [GitHub Discussions](https://github.com/koor-tech/koor/discussions).

Be sure to checkout the [Ceph Common Issues](ceph-common-issues.md) and **that [all prerequisites](../Getting-Started/Prerequisites/prerequisites.md) for the storage backend of your choice** are met.

***

## Where did the `rook-discover-*` Pods go after a recent Rook Ceph update?

A recent change in Rook Ceph has disabled the `rook-discover` DaemonSet by default.
This behavior is controlled by the `ROOK_ENABLE_DISCOVERY_DAEMON` located in the `operator.yaml` or for Helm users `enableDiscoveryDaemon: (false|true` in your values file. It is a boolean, so `false` or `true`.

### When do you want to have `rook-discover-*` Pods / `ROOK_ENABLE_DISCOVERY_DAEMON: true`?

* You are on **(plain) bare metal** and / or simply have "some disks" installed /attached to your server(s), that you want to use for the Rook Ceph cluster.
* If your cloud environment / provider does not provide PVCs with `volumeMode: Block`. Ceph requires block devices (Ceph's `filestore` is not available, through Rook, since a bunch of versions as `bluestore` is superior in certain ways).

## Crash Collector Pods are `Pending` / `ContainerCreating`

* Check the events of the Crash Collector Pod(s) using `kubectl describe pod POD_NAME`.
* If the Pod(s) is waiting for a Secret from the Ceph MONs (keyring for each crash collector), you need to wait a bit longer as the Ceph Cluster is probably still being bootsraped / started up.
* If they are stuck for more than 15-30 minutes, check the Rook Ceph Operator logs if it is stuck in the Ceph Cluster bootstrap / start up procedure.

## No `rook-ceph-mon-*` Pods are running

1. First of all make sure your Kubernetes CNI is working fine! In what feels like 90% of the cases it is network related, e.g., some weird thing with the Kubernetes cluster CNI or other network environment issue.
    * Can you talk to Cluster Service IPs from every node?
    * Can you talk to Pod IPs from every node? Even to Pods not on the same node you are testing from?
    * Check the docs of your CNI, most have a troubleshooting section, e.g., Cilium had some issues from systemd version 245 onwards with `rp_filter`, see here: [rp_filter (default) strict mode breaks certain load balancing cases in kube-proxy-free mode · Issue #13130 · cilium/cilium](https://github.com/cilium/cilium/issues/13130)
2. Does your environment fit all the prerequisites? Check top of page for the links to some of the prerequisites and / or consult the [Koor docs](../Getting-Started/intro.md).
3. Check the `rook-ceph-operator` Logs for any warnings, errors, etc.

## Disk(s) / Partition(s) not used for Ceph

* Does section [When do you want to have `rook-discover-*` Pods / `ROOK_ENABLE_DISCOVERY_DAEMON: true`?](#when-do-you-want-to-have-rook-discover--pods--rook_enable_discovery_daemon-true) apply to you? If so, make sure the operator has the discovery daemon enabled in its (Pod) config!
* Is the disk empty? No leftover partitions on it? Make sure it is either "empty", e.g., nulled by `shred`, `dd` or similar,
    * To make sure the disk is blank as the Rook docs and I recommend the following commands followed by a reboot of the server:
        ```
        DISK="/dev/sdXYZ"
        sgdisk --zap-all "$DISK"
        dd if=/dev/zero of="$DISK" bs=1M count=100 oflag=direct,dsync
        blkdiscard "$DISK"
        ```
        Source: [Koor Storage Cluster Cleanup - Delete the data on hosts](../Storage-Configuration/ceph-teardown.md#delete-the-data-on-hosts)
* Was the disk previously used as a Ceph OSD?
    * Make sure to follow the teardown steps, but make sure to only remove the LVM stuff from that one disk and not from all, see [Koor Storage Cluster Cleanup - Delete the data on hosts](../Storage-Configuration/ceph-teardown.md#delete-the-data-on-hosts).

## A Pod can't mount its PersistentVolume after an "unclean" / "undrained" Node shutdown

1. Check the events of the Pod using `kubectl describe pod POD_NAME`.
2. Check the Node's `dmesg` logs.
3. Check the kubelet logs for errors related to CSI connectivity and / or make sure the node can reach every other Kubernetes cluster node (at least the Rook Ceph cluster nodes (Ceph Mons, OSDs, MGRs, etc.)).
4. Checkout the [CSI Common Issues - Koor Docs](ceph-csi-common-issues.md).

## Ceph CSI: Provisioning, Mounting, Deletion or something doesn't work

Make sure you have checked out the [CSI Common Issues - Koor Docs](ceph-csi-common-issues.md).

If you have some weird kernel and / or kubelet configuration, make sure Ceph CSI's config options in the Rook Ceph Operator config is correctly setup (e.g., `LIB_MODULES_DIR_PATH`, `ROOK_CSI_KUBELET_DIR_PATH`, `AGENT_MOUNTS`).

## Can't run any Ceph Commands in the Toolbox / Ceph Commands timeout

* Are your `rook-ceph-mon-*` Pods all in `Running` state?
* Does a basic `ceph -s` work?
* Is your `rook-ceph-mgr-*` Pod(s) running as well?
* Check the `rook-ceph-mon-*` and `rook-ceph-mgr-*` logs for errors
* Try deleting the toolbox Pod, "maybe it is just a fluke in your Kubernetes cluster network / CNI.
    * Also make sure you are using the latest Rook Ceph Toolbox YAML for the Rook Ceph version you are running on, see [Rook Ceph Toolbox Pod not Creating / Stuck section](#rook-ceph-toolbox-pod-not-creating--stuck).
* In case all these seem to indicate a loss of quorum, e.g., the `rook-ceph-mon-*` talk about `probing` for other mons only, you might need to follow the disaster recovery guide for your Rook Ceph version here: [Rook Ceph Disaster Recovery - Restoring Mon Quorum](disaster-recovery.md#restoring-mon-quorum).

## A MON Pod is running on a Node which is down

* **DO NOT EDIT THE MON DEPLOYMENT!** A MON Deployment can't just be moved to another node without being failovered by the operator and / or if the MON is running using a PVC for its data.
* As long as the operator is running the operator should see the mon being down and fail it over after a configurable timeout.
    * Env var `ROOK_MON_OUT_TIMEOUT`, by default `600s` (10 minutes)

## Remove / Replace a failed disk

Checkout the official Ceph OSD Management guide from Rook here: [Rook Ceph OSD Management Docs](../Storage-Configuration/Advanced/ceph-osd-mgmt.md).

## Rook Ceph Toolbox Pod not Creating / Stuck

* Make sure that you are not using an old version of the Rook Ceph Toolbox, grab the latest manifest here (make sure to switch to the `release-` branch of your Rook release): `https://github.com/rook/rook/blob/master/cluster/examples/kubernetes/ceph/toolbox.yaml`
* The Rook Ceph Toolbox can only fully startup after a Ceph Cluster has at least passed the initial setup by the Rook Ceph operator.
    * Monitor the Rook Ceph Operator logs for errors.
* Check the events of the Toolbox Pod using `kubectl describe pod POD_NAME`.

## Ceph OSD Tree: Wrong Device Class

1. Check device class, second column in `ceph osd tree` output.
2. If you need to change the device class, you first must remove the current one (if it has one set): `ceph osd crush rm-device-class osd.ID`.
3. Now you can set the device class for the OSD: `ceph osd crush set-device-class CLASS osd.ID`
   * Default device classes (at the time of writing): `hdd`, `ssd`, `nvme`
   * Source: [Ceph Docs Latest - CRUSH Maps - Device Classes](https://docs.ceph.com/en/latest/rados/operations/crush-map/#device-classes)

## `HEALTH_WARN: clients are using insecure global_id reclaim` / `HEALTH_WARN: mons are allowing insecure global_id reclaim`

**Source**: https://github.com/rook/rook/issues/7746

> I can confirm this is happening in all clusters, whether a clean install or upgraded cluster, running at least versions: `v14.2.20`, `v15.2.11` or `v16.2.1`.
>
> According to the [CVE also previously mentioned](https://docs.ceph.com/en/latest/security/CVE-2021-20288/), there is a security issue where clients need to be upgraded to the releases mentioned. Once all the clients are updated (e.g. the rook daemons and csi driver), a new setting needs to be applied to the cluster that will disable allowing the insecure mode.
>
> If you see both these health warnings, then either one of the rook or csi daemons has not been upgraded yet, or some other client is detected on the older version:
>
>     health: HEALTH_WARN
>             client is using insecure global_id reclaim
>             mon is allowing insecure global_id reclaim
>
>
> If you only see this one warning, then the insecure mode should be disabled:
>
>     health: HEALTH_WARN
>             mon is allowing insecure global_id reclaim
> To disable the insecure mode from the toolbox after all the clients are upgraded:
> **Make sure all clients have been upgraded, or else those clients will be blocked after this is set**:
>
>     ceph config set mon auth_allow_insecure_global_id_reclaim false
>
> Rook could set this flag automatically after the clients have all been updated.

## Check which "Object Store" is used by an OSD

```console
$ ceph osd metadata 0 | grep osd_objectstore
"osd_objectstore": "bluestore",
```

To get a quick overview of the "object stores" (`bluestore`, (don't use it) `filestore`):
```console
$ ceph osd count-metadata osd_objectstore
{
    "bluestore": 6
}
```

## PersistentVolumeClaims / PersistentVolumes are not Resized

* Make sure the Ceph CSI driver for the storage (block or filesystem) is running (check the logs if you are unsure as well).
* Check if you use a StorageClass that has `allowVolumeExpansion: false`:
    ```console
    $ kubectl get storageclasses.storage.k8s.io
    NAME              PROVISIONER                     RECLAIMPOLICY   VOLUMEBINDINGMODE   ALLOWVOLUMEEXPANSION   AGE
    rook-ceph-block   rook-ceph.rbd.csi.ceph.com      Retain          Immediate           false                  3d21h
    rook-ceph-fs      rook-ceph.cephfs.csi.ceph.com   Retain          Immediate           true                   3d21h
    ```
* To fix this simply set `allowVolumeExpansion: true` in the `StorageClass`. Below is a `StorageClass` with this option set, it is at the top level of the object (not in `.spec` or similar):
    ```yaml hl_lines="1"
    allowVolumeExpansion: true
    apiVersion: storage.k8s.io/v1
    kind: StorageClass
    metadata:
      name: rook-ceph-block
    parameters:
      clusterID: rook-ceph
      csi.storage.k8s.io/controller-expand-secret-name: rook-csi-rbd-provisioner
      [...]
      imageFeatures: layering
      imageFormat: "2"
      pool: replicapool
    provisioner: rook-ceph.rbd.csi.ceph.com
    reclaimPolicy: Retain
    volumeBindingMode: Immediate
    ```

## `[...] failed to retrieve servicemonitor. servicemonitors.monitoring.coreos.com "rook-ceph-mgr" is forbidden: [...]`

You have the Prometheus Operator installed in your Kubernetes cluster, but have not applied the RBAC necessary for the Rook Ceph Operator to be able to create the monitoring objects.

To rectify this, you can run the following command and / or add the file to your deployment system:

```console
kubectl apply -f https://raw.githubusercontent.com/rook/rook/master/cluster/examples/kubernetes/ceph/monitoring/rbac.yaml
```

(Original file located at: [https://github.com/rook/rook/blob/master/cluster/examples/kubernetes/ceph/monitoring/rbac.yaml](https://github.com/rook/rook/blob/master/cluster/examples/kubernetes/ceph/monitoring/rbac.yaml))

## `[...] failed to reconcile cluster "rook-ceph": [...] failed to create servicemonitor. the server could not find the requested resource (post servicemonitors.monitoring.coreos.com)`

This normally means that you don't have the [Prometheus Operator](https://github.com/prometheus-operator/prometheus-operator) installed in your Kubernetes cluster. It is required for `.spec.monitoring.enabled: true` in the CephCluster object to work (the operator to be able to create the `ServiceMonitor` object to enable monitoring).

For the [Rook Ceph Prometheus Monitoring Setup Steps](../Storage-Configuration/Monitoring/ceph-monitoring.md#prometheus-alerts) check the link.

### Solution A: Disable Monitoring in CephCluster

Set `.spec.monitoring.enabled` to `false` in your CephCluster object / yaml (and apply it).

### Solution B: Install Prometheus Operator

If you want to use Prometheus for monitoring your applications and in this case also Rook Ceph Cluster easily in Kubernetes, make sure to install the [Prometheus Operator](https://github.com/prometheus-operator/prometheus-operator).

Checkout the [Prometheus Operator - Getting Started Guide](https://github.com/prometheus-operator/prometheus-operator/blob/main/Documentation/user-guides/getting-started.md).

## `unable to get monitor info from DNS SRV with service name: ceph-mon` / Can't run `ceph` and `rbd` commands in the Rook Ceph XYZ Pod

You are only supposed to run `ceph`, `rbd`, `radosgw-admin`, etc., commands in the **Rook Ceph Toolbox / Tools Pod**.

Regarding the Rook Ceph Toolbox Pod checkout the Rook documentation here: [Ceph Toolbox](ceph-toolbox.md).

### Quick Command to Rook Ceph Toolbox Pod

This requires you to have the Rook Ceph Toolbox deployed, see [Ceph Toolbox](ceph-toolbox.md) for more information.

```console
kubectl -n rook-ceph exec -it $(kubectl -n rook-ceph get pod -l "app=rook-ceph-tools" -o jsonpath='{.items[0].metadata.name}') -- bash
```

## `OSD id X != my id Y` - OSD Crash

1. Exec into a working Ceph OSD on that host `kubectl exec -n rook-ceph -it OSD_POD_NAME -- bash` (`ceph-bluestore-tool` command is needed), run the following commands:
    1. Run `lsblk` to see all disks of the host.
    2. For every disks, run:
        1. Run `ceph-bluestore-tool show-label --dev=/dev/sdX` (note down the OSD ID (`whoami` field in the JSON output) and which disk the OSD is on (example: `OSD 11 /dev/sda`).
2. The `rook-ceph-osd-...` deployment needs to be updated with the new/ correct device path. The `ROOK_BLOCK_PATH` environment variable must have the correct device path (there are two occurrences, in the `containers:` and in `initContainers:` list).
3. After a few seconds / minutes the OSD should show up as `up` in the `ceph osd tree` output (the command can be run in the `rook-ceph-tools` Pod). If you have scaled down the OSD Deployment, make sure to scale it up to `1` again (`kubectl scale -n rook-ceph deployment --replicas=1 rook-ceph-osd...`)

## `_read_bdev_label failed to open /var/lib/ceph/osd/ceph-1/block: (13) Permission denied`

### Issue

* OSD Pod is not starting with logs about the "ceph osd block device" and "permission denied"

### Solution: Do you have the `ceph` package(s) installed on the host and / or a user/group named `ceph`?

This can potentially mess with the owner/group of the ceph osd block device, as described in [GitHub rook/rook Issue 7519 "OSD pod permissions broken, unable to open OSD superblock after node restart"](https://github.com/rook/rook/issues/7519#issuecomment-922263364).

You can either change the user and group ID of the `ceph` user on the host to the one inside the `ceph/ceph` image that your Rook Ceph cluster is running right now (CephCluster object `.spec.cephVersion.image`).

```console
$ kubectl get -n rook-ceph cephclusters.ceph.rook.io rook-ceph -o yaml
apiVersion: ceph.rook.io/v1
kind: CephCluster
metadata:
[...]
  name: rook-ceph
  namespace: rook-ceph
[...]
spec:
  cephVersion:
    image: quay.io/ceph/ceph:v16.2.6-20210927
[...]
```

Depending your hosts, you might not need to even have the `ceph` packages installed. If you are using Rook Ceph, you normally don't need any ceph related packages on the hosts.

Should this have not fixed your issue, you might be running into some other permission issue. If your hosts are using a Linux distribution that uses SELinux, you might need to follow these steps to re-configure the Rook Ceph operator: [OpenShift Special Configuration Guide](../Getting-Started/ceph-openshift.md#rook-settings).

***

Should this page not have yielded you a solution, checkout the [Ceph Common Issues](ceph-common-issues.md) doc as well.

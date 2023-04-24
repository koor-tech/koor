# Ceph Block Storage Backup And Recovery

One of the most crucial parts while installing a rook ceph cluster comes with
having a reliable backup solution, with the ease of restoring those backups in
an efficient manner!

For digging into what would be the best players out there, we experimented with
the existing backup solutions for Rook Ceph.

This series will take you through some of the best viable backup solutions in
hope of avoiding losses and making your lives easier in case of an unfortunate
disaster.

There are several kinds of solutions, but we are going to focus on:

* Ceph features and Ceph native tools
* Kubernetes native tools and projects

## Rook Ceph based Disaster Recovery Solutions

Looking into the first category, based on different workload Disaster Recovery
solutions can narrowed down to:

* Block Storage: RBD Mirroring, which is supported and can be enabled/managed
  for a Block based persistent volume.
* FileSystem Storage: Filesystem Mirroring uses snapshots with one way peering
  for creating a FS based volume backup
* Object Storage: RGW Multisite support

## External Tools based Disaster Recovery Solutions

* [Velero using Restic or Kopia](https://velero.io/docs/v1.3.2/restic/)
* [Kasten.io by Veeam](https://www.kasten.io/)

These can be used to create a snapshot and for the entire application's,
Kubernetes resources. The snapshot and export process can be configured for a
schedule.

Let us dive into each of these solutions and discuss which would suit your use
case best!

## Block Storage Backup (RBD Mirroring)

The Block Storage in Rook Ceph cluster exists in form of Kubernetes's
PersistentVolumes. To backup these resources we'd be using one of Ceph RBDs
feature called [RBD Mirroring](https://docs.ceph.com/en/quincy/rbd/rbd-mirroring/).
RBD Mirroring will asynchronously mirror the RBD images(present in for of
PersistentVolumes) from Primary cluster to a Secondary(Backup) cluster.

This makes use of periodically syncing
[snapshots](https://docs.ceph.com/en/quincy/rbd/rbd-snapshot/), based on `rbd
snap schedule` [read more](https://docs.ceph.com/en/quincy/rbd/rbd-mirroring/).

## Configuring Backup

* On your Rook Ceph Cluster enable RBD mirroring using [Rook’s Block based mirroring official doc](https://rook.io/docs/rook/v1.11/Storage-Configuration/Block-Storage-RBD/rbd-mirroring/#rbd-mirroring)
* Check daemon health status for rbd-mirror

    ```console
    $ kubectl get cephblockpools.ceph.rook.io mirroredpool -n rook-ceph -o jsonpath='{.status.mirroringStatus.summary}'
    {"daemon_health":"OK","health":"OK","image_health":"OK","states":{"replaying":1}}
    ```

* If any issues are found please check the logs and configure the deployment correctly.
* To make sure mirroring is configured correctly, identify the rbd-image mapped to the csi volume using:

    ```console
    # extract rbd_image name
    $ kubectl get pv $&lt;rbd_pv_vol> -o yaml | grep image
    ```

* Go to the toolbox pod and check the mirroring status

    ```console
    $ kubectl exec -it -n rook-ceph deploy/rook-ceph-tools -- bash
    $ rbd mirror image status mirrored-pool/$rbd_image
    ```

    For a health state the rbd mirroring status for the image should look like:

    ```yaml
    test:
      global_id:   05478186-bb6e-4c47-8a3a-da611dabe0a5
      state:       up+stopped
      description: local image is primary
      service:     a on centos-4gb-hel1-7
      last_update: 2022-12-09 09:17:42
      peer_sites:
        name: e7f2c11a-3cea-436c-aa8b-b8b62f93814f
        state: up+replaying
        description: replaying, {"bytes_per_second":0.0,"bytes_per_snapshot":0.0,"local_snapshot_timestamp":1670577394,"remote_snapshot_timestamp":1670577394,"replay_state":"idle"}
        last_update: 2022-12-09 09:17:48
      snapshots:
        9 .mirror.primary.05478186-bb6e-4c47-8a3a-da611dabe0a5.01dfa36d-1963-4536-909a-669be9c1ad51 (peer_uuids:[a692edf 6-bf7c-4bef-aab9-4e4f20530de6])
    ```

When mirroring is enabled and working correctly, you should be able to see the
mirror persistent volume/ rbd image getting synced from primary cluster to
secondary:

On the secondary cluster:

```console
$ kubectl get pv
pvc-b62bf97e-4708-448a-8b04-b5a424034b5f   20Gi       RWO            Retain           Bound      default/rbd-pvc               rook-ceph-block            3d16h
```

Once we have made sure mirroring is working properly, we can create a snap
schedule by running:

```console
$ rbd --cluster site-a mirror snapshot schedule add --pool image-pool 24h 14:00:00-05:00
# check active snap schedule
$ rbd mirror snapshot schedule ls --pool image-pool --recursive
```

Read more about configuring rbd snap schedule [here](https://docs.ceph.com/en/quincy/rbd/rbd-mirroring/#create-image-mirror-snapshots)

Now we have all in place for periodic sync of our block volume in form of rbd
snapshots to a secondary cluster

What if something happens to Block Persistent volume on the primary cluster?

## Restoring Backed up RBD Persistent Volume

If for any reason the primary goes down/something happens to the
PersistentVolume on primary; we can use this **failback process** to restore the
back.

For failback we do the following step by step process using the toolbox pod of
respective clusters:

1. Demote the primary cluster (cluster 1)

    ```console
    cluster_1 $ kubectl exec -it deployments/rook-ceph-tools -- rbd mirror image demote {pool-name}/{image-name}
    ```
    you might observe split brain status momentarily

    ```console
    cluster_1 $ kubectl exec -it deployments/rook-ceph-tools -- rbd mirror image status {pool-name}/{image-name}
    ```

2. Promote the secondary cluster to be the new primary, this step will make the
   cluster 1 (old primary) image to sync cluster 2 (new primary)

    ```console
    cluster_2 $  kubectl exec -it deployments/rook-ceph-tools -- rbd mirror image promote [--force] {pool-name}/{image-name}
    ```

3. And you should be able see rbd image getting synced from cluster2 (new
   primary) to cluster1 (old primary)

    ```console
    cluster_2 $ kubectl exec -it deployments/rook-ceph-tools -- rbd mirror image status {pool-name}/{image-name
    ```

4. Once the sync is complete cluster 2 can be demoted and cluster 1 can be
   promoted back to primary cluster

    ```console
    cluster_2 $ kubectl exec -it deployments/rook-ceph-tools -- rbd mirror image demote {pool-name}/{image-name}
    ```

    ```console
    cluster_1 $  kubectl exec -it deployments/rook-ceph-tools -- rbd mirror image promote {pool-name}/{image-name}
    ```

5. This completes the restore of failed/ corrupt Block Persistent Volume

    ```console
    $ kubectl get pv
    ```
**Note**: The size of the snapshot is expected to be size of the data written in the
image.

For a 3 node Rook Ceph cluster, having 10 GB test data on the Block based PersistentVolume.

| BACKUP STRATEGY | TIME TAKEN TO EXPORT 10GB FILE |
|:---------------:|:------------------------------:|
|  RBD MIRRORING  |             ~31 sec            |

Block Based Backup uses RBD snapshot diff natively for image mirroring, this
will provide higher availability, as secondary cluster can be used when primary
goes down, until restore is finished.

Although, this solution will certainly help with backing up Block Volumes in a
secondary backup cluster in an efficient manner, there might be some [Ceph RBD mirroring concepts](https://docs.ceph.com/en/quincy/rbd/rbd-mirroring/) you
might want to learn with it.

Rook has good documentation around such scenarios as well as planned migration documented [here](https://rook.io/docs/rook/v1.11/Storage-Configuration/Block-Storage-RBD/rbd-async-disaster-recovery-failover-failback/).

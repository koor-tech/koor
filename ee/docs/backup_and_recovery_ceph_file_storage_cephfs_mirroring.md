
# Ceph FileSystem Storage Solution Backup plan

We will enable CephFS mirroring on the primarily-used cluster and add a
secondary (destination) as a peer to it.
The resources from the primary cluster are synchronized to the secondary cluster
at a scheduled time creating a snapshot of the primary filesystem. This results
in regular backups, which can restored whenever needed.

Once a successful peer is added, any volume created on the primary, shall be
synchronized and backed up on the secondary, based on the schedule configured.

Based on your data you can then configure a backup plan.
A backup schedule should incrementally backup the FileSystems resources.
Additionally, we advise that the integrity of the backups to be checked periodically.

## Prerequistes

- You need at least 2x Rook Ceph clusters with at least 3 Nodes (for testing, you can also use [Rook Ceph test cluster manifest](https://github.com/rook/rook/blob/master/deploy/examples/cluster-test.yaml).
The first cluster is used as the source and the second cluster will be used for mirroring CephFS data from the first one.
- For both your Rook Ceph clusters you need to make sure that they can reach each other.
  One of the recommended way to make sure networking is setup correctly is by enabling
  [host networking](https://rook.io/docs/rook/v1.11/Storage-Configuration/Advanced/ceph-configuration/#use-hostnetwork-in-the-cluster-configuration)
  before you deploy the cluster.
- For resource recommendations please refer to [the Rook Documentation](https://rook.io/docs/rook/v1.11/CRDs/Cluster/ceph-cluster-crd/#resource-requirementslimits).

## How to setup?

### Backup

- Have [CephFS Custom Resources Deployed with mirroring enabled](https://rook.io/docs/rook/v1.11/Storage-Configuration/Shared-Filesystem-CephFS/filesystem-mirroring/#create-the-filesystem-with-mirroring-enabled).

- We configured the mirroring schedule to create a backup every 24 hours

``` yaml hl_lines="32 41"
apiVersion: ceph.rook.io/v1
kind: CephFilesystem
metadata:
  name: myfs
  namespace: rook-ceph
spec:
  metadataPool:
    failureDomain: host
    replicated:
      size: 3
  dataPools:
    - name: replicated
      failureDomain: host
      replicated:
        size: 3
  preserveFilesystemOnDelete: true
  metadataServer:
    activeCount: 1
    activeStandby: true
  mirroring:
    enabled: true
    # list of Kubernetes Secrets containing the peer token
    # for more details see: https://docs.ceph.com/en/latest/dev/cephfs-mirroring/#bootstrap-peers
    # Add the secret name if it already exists else specify the empty list here.
    peers:
      secretNames:
        #- secondary-cluster-peer
    # specify the schedule(s) on which snapshots should be taken
    # see the official syntax here https://docs.ceph.com/en/latest/cephfs/snap-schedule/#add-and-remove-schedules
    snapshotSchedules:
      - path: /
        interval: 24h # daily snapshots
        # The startTime should be mentioned in the format YYYY-MM-DDTHH:MM:SS
        # If startTime is not specified, then by default the start time is considered as midnight UTC.
        # see usage here https://docs.ceph.com/en/latest/cephfs/snap-schedule/#usage
        # startTime: 2022-07-15T11:55:00
    # manage retention policies
    # see syntax duration here https://docs.ceph.com/en/latest/cephfs/snap-schedule/#add-and-remove-retention-policies
    snapshotRetention:
      - path: /
        duration: "h 24"
  ```

**Note**: If you are using single node cluster for testing, be sure to change the replicated
size to `1`.

- Deploy CephFS mirror daemon.

```console
cluster1$ kubectl create -f deploy/examples/filesystem-mirror.yaml
```

Once you have the filesystem and mirror daemon deployed, it should look something similar to:

```console
cluster1$ kubectl get pods -A | grep fs
rook-ceph     csi-cephfsplugin-5w6sz                                        2/2     Running     0               41d
rook-ceph     csi-cephfsplugin-btxxs                                        2/2     Running     0               41d
rook-ceph     csi-cephfsplugin-provisioner-75875b5887-k6h8n                 5/5     Running     0               41d
rook-ceph     csi-cephfsplugin-provisioner-75875b5887-w2vqg                 5/5     Running     0               5d1h
rook-ceph     csi-cephfsplugin-q9929                                        2/2     Running     0               41d
rook-ceph     rook-ceph-fs-mirror-7c87686cd7-9dl9n                          2/2     Running     1 (4d15h ago)   4d15h
rook-ceph     rook-ceph-mds-myfs-a-866fc7c5bd-zkw45                         2/2     Running     40 (33h ago)    37d
rook-ceph     rook-ceph-mds-myfs-b-8584dc6bbf-tl7bg                         2/2     Running     42 (33h ago)    37d
```

The steps to configure secondary cluster as a peer can be found [here](https://rook.io/docs/rook/v1.11/Storage-Configuration/Shared-Filesystem-CephFS/filesystem-mirroring/#import-the-token-in-the-destination-cluster).

Please verify the peer setup using the [toolbox pod](https://rook.io/docs/rook/v1.11/Troubleshooting/ceph-toolbox/):

```console
cluster1$ kubectl exec -it -n rook-ceph deploy/rook-ceph-tools -- bash
cluster1$ ceph fs snapshot mirror daemon status | jq
[
  {
    "daemon_id": 6034805,
    "filesystems": [
      {
        "filesystem_id": 1,
        "name": "myfs",
        "directory_count": 3,
        "peers": [
          {
            "uuid": "b92286db-7ee0-40b4-88dd-1ff87d347569",
            "remote": {
              "client_name": "client.mirror",
              "cluster_name": "e7f2c11a-3cea-436c-aa8b-b8b62f93814f",
              "fs_name": "myfs"
            },
            "stats": {
              "failure_count": 0,
              "peer_init_failed": true,
              "recovery_count": 0
            }
          }
        ]
      }
    ]
  }
]
```

Once CephFS peer is set up, you can test the mirroring by creating a test volume:

On the primary cluster,

```console
cluster1$ kubectl exec -n rook-ceph deploy/rook-ceph-tools -t -- ceph fs subvolume create myfs testsubvolume
cluster1$ kubectl exec -n rook-ceph deploy/rook-ceph-tools -t -- ceph fs snapshot mirror enable myfs
cluster1$ kubectl exec -n rook-ceph deploy/rook-ceph-tools -t -- ceph fs snapshot mirror add myfs /volumes/_nogroup/testsubvolume/
```

On the destination (secondary) cluster:

```console
cluster2$ kubectl exec -n rook-ceph deploy/rook-ceph-tools -t -- ceph fs subvolume create myfs testsubvolume
cluster2$ kubectl exec -n rook-ceph deploy/rook-ceph-tools -t -- ceph fs snapshot mirror enable myfs
```

We uploaded a 10 GB test data file to this filesystem and then checked the data integrity of the file on primary cluster, by noting the checksum of the file on the primary, we will use this to verify the file on secondary cluster later on.

```console
# Create the directory
mkdir /tmp/testsubvolume

# Detect the mon endpoints and the user secret for the connection
mon_endpoints=$(grep mon_host /etc/ceph/ceph.conf | awk '{print $3}')
my_secret=$(grep key /etc/ceph/keyring | awk '{print $3}')

# Mount the filesystem
mount -t ceph -o mds_namespace=myfs,name=admin,secret=$my_secret $mon_endpoints:/ /tmp/testsubvolume

# See your mounted filesystem
df -h

$ md5sum testsubvolume/sample.txt
903100bd055574d2bce7b60a054e9751  testsubvolume/sample.txt
```

The snapshot would be created on the scheduled interval but for testing you can create a snapshots and verify if the snap sync is completed.

```console
cluster2$ kubectl exec -n rook-ceph deploy/rook-ceph-tools -t -- ceph fs subvolume snapshot create myfs testsubvolume snap3
```

On the secondary, we verify if these snapshots got synced.

```console
cluster2$ ceph --admin-daemon /var/run/ceph/ceph-client.fs-mirror.11.94285354653568.asok fs mirror peer status myfs@1 64b01275-ee93-4b5b-a9d4-bda551ef4db0
{
    "/volumes/_nogroup/testsubvolume": {
        "state": "idle",
        "last_synced_snap": {
            "id": 5,
            "name": "snap3",
            "sync_duration": 0.12599942,
            "sync_time_stamp": "752855.084361s"
        },
        "snaps_synced": 1,
        "snaps_deleted": 0,
        "snaps_renamed": 0
    },
    "/volumes/_nogroup/testsubvolume2": {
        "state": "idle",
        "last_synced_snap": {
            "id": 9,
            "name": "snap3",
            "sync_duration": 0.11599946799999999,
            "sync_time_stamp": "755966.611134s"
        },
        "snaps_synced": 1,
        "snaps_deleted": 0,
        "snaps_renamed": 0
    }
}
```

We see, snap sync count got incremented, we will be able to see a persistent volume got created on the secondary cluster.

Now we verify on the secondary(destination) cluster, the checksum of the file after volume mount we verify if the snapshot after file being written got synchronized...

```console
# Enter into toolbox or direct mount pod
cluster2$ kubectl exec -it -n rook-ceph deploy/rook-ceph-tools -- bash

# Create the directory
mkdir /tmp/testsubvolume

# Detect the mon endpoints and the user secret for the connection
mon_endpoints=$(grep mon_host /etc/ceph/ceph.conf | awk '{print $3}')
my_secret=$(grep key /etc/ceph/keyring | awk '{print $3}')

# Mount the filesystem
mount -t ceph -o mds_namespace=myfs,name=admin,secret=$my_secret $mon_endpoints:/ /tmp/testsubvolume

# See your mounted filesystem
df -h

# verify checksum of the file synchronized with filesystem
$ md5sum testsubvolume/sample.txt
903100bd055574d2bce7b60a054e9751  testsubvolume/sample.txt
```

To verify the data, we followed the steps to mount the volume using direct mount pod, which
can be found in detail on the Rook Documentation [here](https://rook.io/docs/rook/v1.11/Troubleshooting/direct-tools/#shared-filesystem-tools).

Voilà! We have successfully configured backup for Rook Ceph Persistent Volumes using CephFS mirroring.

### Recovery

In case of any failure (disaster), since CephFS mirroring is a one way mirroring you would
need to remove the existing peer, and reverse the roles i.e. make the destination cluster new primary and primary, the destination, so that we can sync back the snapshots.
This sounds tedious but is a simple process that can be used at the time of restoration.

From the [toolbox pod](https://rook.io/docs/rook/v1.11/Troubleshooting/ceph-toolbox/) run the following commands.

- Unlink the disaster struck primary (cluster1) peer from the secondary cluster.

```console
cluster2$ ceph fs snapshot mirror peer_remove myfs <peer_uuid>
# verify from peer list
cluster2$ ceph fs snapshot mirror peer_list myfs
```

- Bring back the disaster struck primary (cluster1) to healthy state and make it secondary by adding it as a peer for
  the current primary (cluster2).

```console
cluster2$ ceph fs snapshot mirror peer_bootstrap create myfs client.admin
{"token": "eyJmc2lkIjogIjBkZjE3MjE3LWRmY2QtNDAzMC05MDc5LTM2Nzk4NTVkNDJlZiIsICJmaWxlc3lzdGVtIjogImJhY2t1cF9mcyIsICJ1c2VyIjogImNsaWVudC5taXJyb3JfcGVlcl9ib290c3RyYXAiLCAic2l0ZV9uYW1lIjogInNpdGUtcmVtb3RlIiwgImtleSI6ICJBUUFhcDBCZ0xtRmpOeEFBVnNyZXozai9YYUV0T2UrbUJEZlJDZz09IiwgIm1vbl9ob3N0IjogIlt2MjoxOTIuMTY4LjAuNTo0MDkxOCx2MToxOTIuMTY4LjAuNTo0MDkxOV0ifQ=="}
```

```console
cluster1$ ceph fs snapshot mirror peer_bootstrap import myfs eyJmc2lkIjogIjBkZjE3MjE3LWRmY2QtNDAzMC05MDc5LTM2Nzk4NTVkNDJlZiIsICJmaWxlc3lzdGVtIjogImJhY2t1cF9mcyIsICJ1c2VyIjogImNsaWVudC5taXJyb3JfcGVlcl9ib290c3RyYXAiLCAic2l0ZV9uYW1lIjogInNpdGUtcmVtb3RlIiwgImtleSI6ICJBUUFhcDBCZ0xtRmpOeEFBV      nNyZXozai9YYUV0T2UrbUJEZlJDZz09IiwgIm1vbl9ob3N0IjogIlt2MjoxOTIuMTY4LjAuNTo0MDkxOCx2MToxOTIuMTY4LjAuNTo0MDkxOV0i
```

- Verify peer on cluster2(current primary) that cluster1 was successfully added
  as the peer.

```console
cluster1$ ceph fs snapshot mirror peer_list myfs
```

- After the peer is configured, the snapshots should be synchronized back to the previous primary(cluster1), you should be able to see Persistent Volumes created on cluster1.

- Once you are sure all the resources are synchronized back to cluster1(current
  secondary), you can remove the peer and make current secondary(cluster1) back
  the primary and cluster2 can be peered back to resume backup on them. Please
  ensure, the recovery be performed ensuring the primary cluster where we are
  restoring backups is in healthy state.

```console
cluster1$ kubectl exec -n rook-ceph deploy/ceph-fs-mirror -- bash
cluster1$ ls /var/run/ceph/
ceph-client.fs-mirror.11.94285354653568.asok
cluster1$  ceph --admin-daemon /var/run/ceph/ceph-client.fs-mirror.11.94285354653568.asok fs mirror peer status myfs@1 64b01275-ee93-4b5b-a9d4-bda551ef4db0
    "/volumes/_nogroup/testsubvolume": {
        "state": "idle",
        "last_synced_snap": {
            "id": 26,
            "name": "snap3",
            "sync_duration": 89.813304341000006,
            "sync_time_stamp": "772137.092616s"
        },
        "snaps_synced": 1,
        "snaps_deleted": 0,
        "snaps_renamed": 0
```

## Results

We ran Rook Ceph on hardware, with 16GB of RAM on each of the 3 nodes.
The result was an effortless and successful synchronization of the CephFS Persistent
Volumes with the sample data integrity intact, it took around 89 sec for a
snapshot to be synchronized.

**Note**: The results may vary depending on the hardware used and network bandwidth of the cluster, along with other factors.


# Glossary

## Ceph

Ceph is a distributed network storage and file system with distributed metadata management and POSIX semantics.

### MON

Ceph Monitor is a daemon that maintains a map of the state of the cluster. This “cluster state” includes the monitor map, the manager map, the OSD map, and the CRUSH map. A Ceph cluster must contain a minimum of three running monitors in order to be both redundant and highly-available. Ceph monitors and the nodes on which they run are often referred to as “mon”s.

### OSD

Ceph Object Storage Daemon. The Ceph OSD software, which interacts with logical disks.

### MGR 

The Ceph manager software, which collects all the state from the whole cluster in one place.

### RADOS

Reliable Autonomic Distributed Object Store. RADOS is the object store that provides a scalable service for variably-sized objects. The RADOS object store is the core component of a Ceph cluster.

### RADOS Cluster

A proper subset of the Ceph Cluster consisting of OSDs, Ceph Monitors, and Ceph Managers.

### MDS

The Ceph MetaData Server daemon. Also referred to as “ceph-mds”. The Ceph metadata server daemon must be running in any Ceph cluster that runs the CephFS file system. The MDS stores all filesystem metadata.

### RBD

The block storage component of Ceph. Also called "RADOS Block Device". Its a software instrument that orchestrates the storage of block-based data in Ceph.

### RGW

RADOS Gate Way. The component of Ceph that provides a gateway to both the Amazon S3 RESTful API and the OpenStack Swift API.

For more Ceph related terms, please refer to the [Ceph Glossary page](https://docs.ceph.com/en/latest/glossary/).

## Kubernetes

Kubernetes, also known as K8s, is an open-source system for automating deployment, scaling, and management of containerized applications.

### OpenShift

OpenShift Container Platform is a cloud-based Kubernetes container platform. The foundation of OpenShift Container Platform is based on Kubernetes and therefore shares the same technology.

### Object Bucket Claim (OBC)

An Object Bucket Claim (OBC) is custom resource which requests a bucket (new or existing).

### Object Bucket (OB)

An Object Bucket (OB) is a custom resource automatically generated when a bucket is provisioned. It is a global resource, typically not visible to non-admin users, and contains information specific to the bucket.
    
### nodeSelector

nodeSelector is the simplest recommended form of node selection constraint. You can add the nodeSelector field to your Pod specification and specify the node labels you want the target node to have. Kubernetes only schedules the Pod onto nodes that have each of the labels you specify. For further information please refer to [official Kubernetes documentation](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#nodeselector).

### Node affinity

Node affinity is conceptually similar to nodeSelector, allowing you to constrain which nodes your Pod can be scheduled on based on node labels.

### CSI

The Container Storage Interface (CSI) is a standard for exposing arbitrary block and file storage systems to containerized workloads on Container Orchestration Systems (COs) like Kubernetes. For further information please refer to [official CSI documentation](https://kubernetes-csi.github.io/docs/introduction.html).

### Storage Classes

A [StorageClass](https://kubernetes.io/docs/concepts/storage/storage-classes/#introduction) provides a way for administrators to describe the "classes" of storage they offer.

### PV

A PersistentVolume (PV) is a piece of storage in the cluster that has been provisioned by an administrator or dynamically provisioned using Storage Classes. It is a resource in the cluster just like a node is a cluster resource.

### PVC

A PersistentVolumeClaim (PVC) is a request for storage by a user. It is similar to a Pod. Pods consume node resources and PVCs consume PV resources. Pods can request specific levels of resources (CPU and Memory). Claims can request specific size and access modes (e.g., they can be mounted ReadWriteOnce, ReadOnlyMany or ReadWriteMany). Please refer [Access Modes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes)

### CustomResourceDefinitions or CRDs

The [CustomResourceDefinition](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/#customresourcedefinitions) API resource allows you to define custom resources. Defining a CRD object creates a new custom resource with a name and schema that you specify. The Kubernetes API serves and handles the storage of your custom resource.

### Finalizers

[Finalizers](https://kubernetes.io/docs/concepts/overview/working-with-objects/finalizers/) are namespaced keys that tell Kubernetes to wait until specific conditions are met before it fully deletes resources marked for deletion. Finalizers alert controllers to clean up resources the deleted object owned.
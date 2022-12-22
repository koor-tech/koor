---
title: "Troubleshooting"
---

To help troubleshoot your Koor Storage Distribution clusters, here are some tips on what information will help solve the issues you might be seeing.
If after trying the suggestions found on this page and the problem is not resolved, the Koor Storage Distribution team is very happy to help you troubleshoot the issues through [GitHub Discussions](https://github.com/koor-tech/koor/discussions).

## Troubleshooting Techniques

Kubernetes status and logs are the main resources needed to investigate issues in any Rook cluster.

## Kubernetes Tools

Kubernetes status is the first line of investigating when something goes wrong with the cluster. Here are a few artifacts that are helpful to gather:

### Rook Pod Status

```console
kubectl get -n <cluster-namespace> pod -o wide
kubectl get -n rook-ceph pod -o wide
```

### Logs for Rook Pods

Logs for the operator Pod:

```console
kubectl logs -n <cluster-namespace> -l app=<storage-backend-operator>
kubectl logs -n rook-ceph -l app=rook-ceph-operator
```

Logs for a specific pod:

```console
kubectl logs -n <cluster-namespace> <pod-name>
```

Logs of a pod selected using labels, such as `mon=a`:

```console
kubectl logs -n <cluster-namespace> -l <label-matcher>
kubectl logs -n rook-ceph -l mon=a
```

Some pods have specialized init containers, so you may need to look at logs for different containers
within the pod. To get the logs for a specific container of a pod:

```console
kubectl logs -n <namespace> <pod-name> -c <container-name>
kubectl logs -n rook-ceph rook-ceph-mon-b-55549bc497-phgtz -c chown-container-data-dir
```

#### Pods with multiple containers

For all container logs:

```console
kubectl logs -n <cluster-namespace> <pod-name> --all-containers
kubectl logs -n rook-ceph rook-ceph-mon-b-55549bc497-phgtz --all-containers
```

For a single container:

```console
kubectl logs -n <cluster-namespace> <pod-name> -c <container-name>
kubectl logs -n rook-ceph rook-ceph-mon-b-55549bc497-phgtz -c mon
```

Logs for pods which are no longer running:

```console
kubectl logs -n <cluster-namespace> <pod-name> --previous
kubectl logs -n rook-ceph rook-ceph-mon-b-55549bc497-phgtz --previous
```

### Logs from a specific Node/ Server

To find why a PVC is failing to mount:

1. Check the Pod Events (at the bottom of the command's output)

    ```console
    kubectl describe pod <pod-name>
    ```

2. Verify that the Node is shown as `Ready` in your cluster using:

    ```console
    kubectl get nodes
    ```

3. Connect to the node (e.g., use `ssh`)
4. Get the kubelet service logs (if your distro is using systemd)

    ```console
    journalctl -u kubelet
    ```

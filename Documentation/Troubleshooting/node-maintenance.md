---
title: Node Maintenance
---

## How to perform node maintenance for any of the nodes?

Here are the steps for the same:

1. Get the node and pod details :

    ```bash
    kubectl get pods -owide -n rook-ceph |egrep '<node-name>'
    ```

2. Make the node unschedulable :

    ```bash
    kubectl cordon <nodename>
    ```

3. Drain the node :

    ```bash
    kubectl drain <node name>
    ```

    For more details, please refer to [Kubernetes documentation](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#drain).

4. Shut down / Reboot the node for maintenance.

5. Once the node is back up uncordon the node :

    ```bash
    kubectl uncordon <node name>
    ```

### **NOTE**

For OSD Maintenance and management, please refer to [Ceph OSD Management](https://rook.io/docs/rook/latest/Storage-Configuration/Advanced/ceph-osd-mgmt/) from [rook.io](https://rook.io/) documentation.

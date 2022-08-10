# Roadmap

This document defines a high level roadmap for the Koor Storage Distribution development and upcoming releases.
The features and themes included in each milestone are optimistic in the sense that some do not have clear owners yet.
Community and contributor involvement is vital for successfully implementing all desired items for each release.
We hope that the items listed below will inspire further engagement from the community to keep Koor Storage Distribution progressing and shipping exciting and valuable features.

Any dates listed below and the specific issues that will ship in a given milestone are subject to change but should give a general idea of what we are planning.
See the [GitHub project boards](https://github.com/koor-tech/koor/projects) for the most up-to-date issues and their status.

## Currently

### Basic Setup: Naming and building on the Rook codebase (2-3 weeks by mid-August)

* This is about getting all CI workflows running (e.g., testing, building and release flows).
* Initial code optimization are included with this.

## Upcoming

### Improved Cluster Monitoring (3-4 weeks by mid-September)

* We will create a separate project for a Prometheus-based exporter.
* The exporter will be build modularly to make it easy to expand it in the future.
* The initial module created for the exporter will gather metrics about object storage usage/ quotas. This is to improve visibility into application's object storage usage.
* For every metrics module, we will look into updating existing and/ or creating new Grafana dashboards for visualization of these metrics.

### Easier Ceph Dashboard SSO Setup (1-2 weeks by mid-September, concurrently)

* Phase 1: Providing re-built Ceph images that include the necessary SSO libraries (OAuth2, OpenID, etc.)
* Phase 2: Depending on our customer's feedback, we are going to look into a Kubernetes native way to configure SSO on the Ceph MGR dashboard (e.g., expansion of existing Custom Resource Objects).

### First Long-Term Supported Stable Version of the Koor Storage Distribution (4-6 weeks, by mid-November)

* We still need to iron out the details for a long-term supported version at this moment (e.g., time span, how to handle code improvements).
* We are going to discuss with the Rook community, on how we can approach this together and work out a plan from there.

### Backup/ Restore Flows/ Processes

* Phase 1: Integration with existing backup & restore projects/products (3-4 weeks, end of year, continuously being worked on)
* Phase 2: Our own simple but effective solution to take backups and restore them (initial approach 3-4 weeks, more flexible/dynamic backups in other environments over longer time)

### Optimized Cluster Cloud Topology Handling (5-6 weeks, at latest mid January)

* Improvements to the operator logic in handling cluster node changes (e.g. removal/ addition of new storage nodes).
* This will make it easier to handle new nodes, especially in Cloud environments with regions, availability zones/ sections, to ensure an improved storage availability.

***

> **Disclaimer**: Our development roadmap is subject to change, do not base your purchase on this.

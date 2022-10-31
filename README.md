# Koor Storage Distribution

[![GitHub release](https://img.shields.io/github/release/koor-tech/koor/all.svg)](https://github.com/koor-tech/koor/releases)
[![Docker Pulls](https://img.shields.io/docker/pulls/koor-tech/ceph)](https://hub.docker.com/u/koorinc)
[![Go Report Card](https://goreportcard.com/badge/github.com/koor-tech/koor)](https://goreportcard.com/report/github.com/koor-tech/koor)
[![Security scanning](https://github.com/koor-tech/koor/actions/workflows/synk.yaml/badge.svg)](https://github.com/koor-tech/koor/actions/workflows/synk.yaml)
[![Twitter Follow](https://img.shields.io/twitter/follow/koor_tech.svg?style=social&label=Follow)](https://twitter.com/intent/follow?screen_name=koor_tech&user_id=1509666502714265604)

## What is Koor Storage Distribution?

The Koor Storage Distribution is an open source **cloud-native storage orchestrator** for Ceph on Kubernetes, forked from the [Rook open source project](https://github.com/rook/rook).

Koor Storage Distribution turns Ceph into self-managing, self-scaling, and self-healing storage services. It does this by automating deployment, bootstrapping, configuration, provisioning, scaling, upgrading, migration, disaster recovery, monitoring, and resource management. Koor Storage Distribution uses the facilities provided by the underlying cloud-native container management, scheduling and orchestration platform to perform its duties.

Koor Storage Distribution integrates deeply into cloud native environments leveraging extension points and providing a seamless experience for scheduling, lifecycle management, resource management, security, monitoring, and user experience.

The Rook project is hosted by the [Cloud Native Computing Foundation](https://cncf.io) (CNCF) as a [graduated](https://www.cncf.io/announcements/2020/10/07/cloud-native-computing-foundation-announces-rook-graduation/) level project. If you are a company that wants to help shape the evolution of technologies that are container-packaged, dynamically-scheduled and microservices-oriented, consider joining the CNCF. For details about who's involved and how Rook plays a role, read the CNCF [announcement](https://www.cncf.io/blog/2018/01/29/cncf-host-rook-project-cloud-native-storage-capabilities).

## Getting Started and Documentation

For installation, deployment, and administration, see our [Documentation](https://docs.koor.tech/docs/latest).

## Contributing

We welcome contributions. See [Contributing](CONTRIBUTING.md) to get started.

## Report a Bug

For filing bugs, suggesting improvements, or requesting new features, please open an [issue](https://github.com/koor-tech/koor/issues).

### Reporting Security Vulnerabilities

If you find a vulnerability or a potential vulnerability in Koor Storage Distribution please let us know immediately at
[security@koor.tech](mailto:security@koor.tech). We'll send a confirmation email to acknowledge your
report, and we'll send an additional email when we've identified the issues positively or
negatively.

For further details, please see the complete [security release process](SECURITY.md).

## Contact

Please use the following to reach members of the community:

- GitHub: Start a [discussion](https://github.com/koor-tech/koor/discussions) or open an [issue](https://github.com/koor-tech/koor/issues)
- Twitter: [@koor_tech](https://twitter.com/koor_tech)
- Security topics: [security@koor.tech](#reporting-security-vulnerabilities)
- Office hours meeting: We have a bi-weekly office hour meeting to answer questions and help with any issues around Koor, Rook and Ceph, see [office hours section](#office-hours).

## Office Hours

This is our bi-weekly office hour meeting to answer any questions and help with issues around Koor, Rook and Ceph.

You can add the meeting to your calendar using [this link (invite is sent through Google Calendar)](https://calendar.google.com/calendar/event?action=TEMPLATE&tmeid=NHRhMTBqY2Y0ZTFkb2x1MnZkYThma290M2FfMjAyMjExMDlUMTgzMDAwWiBjXzJjY2Y0OWY1NDZlYzRlYzQ0NzhhMmRiMDI1ZmVjYjdmN2U4MDgxMjZkYmViNzY3MWYxMzg1NGVlNjgwNmQyMmRAZw&tmsrc=c_2ccf49f546ec4ec4478a2db025fecb7f7e808126dbeb7671f13854ee6806d22d%40group.calendar.google.com&scp=ALL).

Feel free to add a topic to the meeting minutes document [here](https://docs.google.com/document/d/1twakYk3XNZD_1Xmi3GDXojuPPkUp7fb06e_4rtgNWdM/edit).

## Official Releases

Official releases of Koor Storage Distribution can be found on the [releases page](https://github.com/koor-tech/koor/releases).
Please note that it is **strongly recommended** that you use [official releases](https://github.com/koor-tech/koor/releases) of Koor Storage Distribution, as unreleased versions from the master branch are subject to changes and incompatibilities that will not be supported in the official releases.
Builds from the master branch can have functionality changed and even removed at any time without compatibility support and without prior notice.

## Licensing

See the [LICENSE](LICENSE) file for licensing information for this repository.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fkoor-tech%2koor.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fkoor-tech%2koor?ref=badge_large)

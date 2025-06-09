# notification-controller

[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/4791/badge)](https://bestpractices.coreinfrastructure.org/projects/4791)
[![e2e](https://github.com/fluxcd/notification-controller/workflows/e2e/badge.svg)](https://github.com/fluxcd/notification-controller/actions)
[![report](https://goreportcard.com/badge/github.com/fluxcd/notification-controller)](https://goreportcard.com/report/github.com/fluxcd/notification-controller)
[![license](https://img.shields.io/github/license/fluxcd/notification-controller.svg)](https://github.com/fluxcd/notification-controller/blob/main/LICENSE)
[![release](https://img.shields.io/github/release/fluxcd/notification-controller/all.svg)](https://github.com/fluxcd/notification-controller/releases)

## VGS fork

```
today=$(date +%Y%m%d)
export IMG="190066226418.dkr.ecr.us-east-1.amazonaws.com/infra/notification-controller:v0.31.0-$today" 
make docker-build
assume presidio/infra
make docker-push
```

Event forwarder and notification dispatcher for the [GitOps Toolkit](https://fluxcd.io/flux/components/) controllers.
The notification-controller is an implementation of the [notification.toolkit.fluxcd.io](docs/spec/v1beta1/README.md)
API based on the specifications described in the [RFC](docs/spec/README.md).

![overview](docs/diagrams/notification-controller-overview.png)

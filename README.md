Supergiant: Kubernetes Orchestration
=


<!-- Links -->

[Kubernetes Source URL]: https://github.com/kubernetes/kubernetes
[Supergiant Website URL]: https://supergiant.io/
<!-- [Supergiant Docs URL]: https://supergiant.io/docs -->
[Supergiant Tutorials URL]: https://supergiant.io/tutorials
[Supergiant Slack URL]: https://supergiant.io/slack
[Supergiant Community URL]: https://supergiant.io/community
[Supergiant Contribution Guidelines URL]: http://supergiant.github.io/docs/community/contribution-guidelines.html
<!-- [Supergiant Swagger Docs URL]: http://swagger.supergiant.io/docs/ -->
[Tutorial AWS URL]: https://supergiant.io/blog/how-to-install-supergiant-container-orchestration-engine-on-aws-ec2?utm_source=github
[Tutorial Linux URL]: https://supergiant.io/blog/how-to-start-supergiant-server-as-a-service-on-ubuntu?utm_source=github
[Tutorial MongoDB URL]: https://supergiant.io/blog/deploy-a-mongodb-replica-set-with-docker-and-supergiant?urm_source=github
[Community and Contributing Anchor]: #community-and-contributing
<!-- [Swagger URL]: http://swagger.io/ -->
[Git URL]: https://git-scm.com/
[Go URL]: https://golang.org/
[Go Remote Packages URL]: https://golang.org/doc/code.html#remote
[Supergiant Go Package Anchor]: #how-to-install-supergiant-as-a-go-package
[Generate CSR Anchor]: #how-to-generate-a-certificate-signing-request-file
<!-- [Create Admin User Anchor]: #create-an-admin-user -->
[Install Dependencies Anchor]: #installing-generating-dependencies

<!-- Badges -->

[GoReportCard Widget]: https://goreportcard.com/badge/github.com/supergiant/supergiant
[GoReportCard URL]: https://goreportcard.com/report/github.com/supergiant/supergiant
[GoDoc Widget]: https://godoc.org/github.com/supergiant/supergiant?status.svg
[GoDoc URL]: https://godoc.org/github.com/supergiant/supergiant
[Govendor URL]: https://github.com/kardianos/govendor
[Travis Widget]: https://travis-ci.org/supergiant/supergiant.svg?branch=master
[Travis URL]: https://travis-ci.org/supergiant/supergiant
[Release Widget]: https://img.shields.io/github/release/supergiant/supergiant.svg
[Release URL]: https://github.com/supergiant/supergiant/releases/latest
[Coverage Status]: https://coveralls.io/github/supergiant/supergiant?branch=master
[Coverage Status Widget]: https://coveralls.io/repos/github/supergiant/supergiant/badge.svg?branch=master
<!-- [Swagger API Widget]: http://online.swagger.io/validator?url=http://swagger.supergiant.io/api-docs -->
<!-- [Swagger URL]: http://swagger.supergiant.io/docs/ -->

# <img src="http://supergiant.io/img/logo_dark.svg" width="400">

[![GoReportCard Widget]][GoReportCard URL] [![GoDoc Widget]][GoDoc URL] [![Release Widget]][Release URL] [![Travis Widget]][Travis URL] [![Coverage Status Widget]][Coverage Status]

---

Supergiant empowers developers and administrators through its simplified deployment and management of Kubernetes, in addition to easing the configuration and deployment of Helm charts, taking advantage of Kubernetes' power, flexibility, and abstraction.

Supergiant facilitates clusters on multiple cloud providers, striving for truly agnostic and impartial infrastructure--and it does this with an autoscaling system that cares deeply about efficiency. It asserts through downscaling and resource packing that unutilized infrastructure shouldn't be paid for (and, therefore, shouldn't be running).

Supergiant implements simple practices that abstract load-balancing, application deployment, basic monitoring, node deployment or destruction, and more, on a highly usable UI. Its efficient packing algorithm enables seamless auto-scaling of Kubernetes clusters, minimizing costs while maintaining the resiliency of applications. To dive into top-level concepts, see [the documentation](https://supergiant.readthedocs.io/en/docs/API/capacity_service/).

# Features

* Fully compatible with native Kubernetes versions 1.5.7, 1.6.7, 1.7.7, and 1.8.7
* Easy management and deployment of multiple kubes in various configurations
* AWS, DigitalOcean, OpenStack, Packet, GCE, and on-premise kube deployment
* Easy creation of Helm releases, Pods, Services, LoadBalancers, etc.
* Automatic, resource-based node scaling
* Compatibility with multiple hardware architectures
* Role-based Users, Session-based logins, self-signed SSLs, and API tokens
* A clean UI and CLI, both built on top of an API (with importable [Go client lib](pkg/client))

# Micro-Roadmap

Currently, the core team is working on the following:

* Add LDAP and OAuth user authentication
* Add support for new cloud providers
* Add support for local installations

# Resources

- [Supergiant Website][Supergiant Website URL]
- [Top-level concepts](https://supergiant.readthedocs.io/en/docs/API/capacity_service/))
- [Tutorials](https://supergiant.io/tutorials)
- [Slack Support Channel](https://supergiant.io/slack)
- [Installation](https://supergiant.readthedocs.io/en/docs/Installation/Linux/)
~ [UI Usage](http://supergiant.readthedocs.io/en/docs/Using%20the%20UI/cloud_accounts/)
~ [API Usage](http://supergiant.readthedocs.io/en/docs/Using%20the%20API/load_balancer/)

# Community and Contributing

We are grateful for any contribution to the Supergiant project be it in a form of a new GitHub issue, a GitHub feature Pull Request, social media engagement etc. Contributing to Supergiant projects requires familiarization with Community and our Contribution Guidelines. Please see these links to get started.

* [Community Page][Supergiant Community URL]
* [Contribution Guidelines][Supergiant Contribution Guidelines URL]


# Development

## Use Docker in development

    docker-compose build server
    docker-compose run --rm --service-ports server

## Native go on your host

If you would like to contribute changes to Supergiant, first see the pages in
the section above, [Community and Contributing][Community and Contributing Anchor].

_Note: [Supergiant cloud installers][Tutorial AWS URL] have dependencies
pre-installed and configured and will generate a self-signed cert based on the
server hostname. These instructions are for setting up a local or custom
environment._

Supergiant dependencies:

* [Git][Git URL]
* [Go][Go URL] version 1.7+
* [Govendor][Govendor URL] for vendoring Go dependencies

## Checkout the repo

```shell
go get github.com/supergiant/supergiant
```

## Create a Config file

You can copy the [example configuration](config/config.json.example):

```shell
cp config/config.json.example config/config.json
```

## Run Supergiant

```shell
go run cmd/server/server.go --config-file config/config.json
open localhost:8080
```

## Build the CLI

This will allow for calling the CLI with the `supergiant` command:

```shell
go build -o $GOPATH/bin/supergiant cmd/cli/cli.go
```

## Run Tests

```shell
govendor test +local
```

## Saving dependencies

If you make a change and import a new package, run this to vendor the imports.

```shell
govendor add +external
```

## Compiling Provider files, UI templates, and static assets

Supergiant uses [go-bindata](https://github.com/jteeuwen/go-bindata) to compile
assets directly into the code. You will need to run this command if you're
making changes to the UI _or_ if you're working with Provider code:

```shell
go-bindata -pkg bindata -o bindata/bindata.go config/providers/... ui/assets/... ui/views/...
```

## Enabling SSL

Our AMI distribution automatically sets up self-signed SSL for Supergiant, but
the default [config/config.json.example](config/config.json.example)
does not enable SSL.

You can see [our AMI boot file](build/sgboot) for an example of how
that is done if you would like to use SSL locally or on your own production
setup.

# License

This software is licensed under the Apache License, version 2 ("ALv2"), quoted below.

Copyright 2016 Qbox, Inc., a Delaware corporation. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License"); you may not
use this file except in compliance with the License. You may obtain a copy of
the License at http://www.apache.org/licenses/LICENSE-2.0.

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

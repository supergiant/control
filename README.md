SUPERGIANT: Easy container orchestration using Kubernetes
=========================================================

---

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

### <img src="http://supergiant.io/img/logo_dark.svg" width="400">

[![GoReportCard Widget]][GoReportCard URL] [![GoDoc Widget]][GoDoc URL] [![Travis Widget]][Travis URL] [![Release Widget]][Release URL] [![Coverage Status Widget]][Coverage Status]

---


Supergiant is an open-source container orchestration system that lets developers
easily deploy and manage apps as Docker containers using Kubernetes.

We want to make Supergiant the easiest way to run Kubernetes in the cloud.

Quick start...

* [How to Install Supergiant Container Orchestration Engine on AWS EC2][Tutorial AWS URL]
<!-- * [Deploy a MongoDB Replica Set with Docker and Supergiant][Tutorial MongoDB URL]
  _(Note: this tutorial is out of date, but you can see
  [the current example here.](examples/deploy_mongo.sh))_ -->

---

![Supergiant UI](http://g.recordit.co/EfUk4D863W.gif)

## Features

* Fully compatible with native Kubernetes (works with existing setups)
* UI and CLI, both built on top of an API (with importable
  [Go client lib](pkg/client))
* Filterable container metrics views (RAM / CPU timeseries graphs)
* Deploy / Update / Restart containers with a few clicks
* Launch and manage multiple Kubes across multiple cloud providers from the UI
* Works with multiple cloud providers (AWS, DigitalOcean, OpenStack, and
  _actively_ adding more, in addition to on-premise hardware support)
* Automatic server management (background server autoscaling, up/down depending
  on container resource needs)
* Role-based Users, Session-based login, self-signed SSL, and API tokens for
  security (OAuth and LDAP support to come)



## Resources

* [Supergiant Website][Supergiant Website URL]
* [Docs](docs/v0/)
* [Tutorials](https://supergiant.io/tutorials)
* [Slack](https://supergiant.io/slack)
* [Install][Tutorial AWS URL]


## Installation

Checkout the [releases](https://github.com/supergiant/supergiant/releases) page
for binaries on Windows, Mac, and Linux.

Copy (and customize if necessary)
[the example config file](config/config.json.example), and run with:

```shell
<supergiant-server-binary> --config-file config.json
```

If you want to easily install Supergiant on Amazon Web Services EC2, follow the
[Supergiant Install Tutorial][Tutorial AWS URL].


## Top-Level Concepts

See [the docs folder](docs/v0/).

<!-- Supergiant makes use of the [Swagger API framework][Swagger URL] for documenting
all resources. See the full Supergiant API documentation for the full reference.

* [![Swagger API Widget] Supergiant Swagger API reference][Supergiant Swagger
Docs URL] -->


## Micro-Roadmap

Currently, the core team is working on the following:

* Add LDAP and OAUTH user authentication
* Add support for additional cloud providers
* Add support for local installations


## Community and Contributing

We are very grateful of any contribution.

All Supergiant projects require familiarization with our Community and our
Contribution Guidelines. Please see these links to get started.

* [Community Page][Supergiant Community URL]
* [Contribution Guidelines][Supergiant Contribution Guidelines URL]


## Development

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

#### Checkout the repo

```shell
go get github.com/supergiant/supergiant
```

#### Create a Config file

You can copy the [example configuration](config/config.json.example):

```shell
cp config/config.json.example config/config.json
```

#### Run Supergiant

```shell
go run cmd/server/server.go --config-file config/config.json
open localhost:8080
```

#### Build the CLI

This will allow for calling the CLI with the `supergiant` command:

```shell
go build -o $GOPATH/bin/supergiant cmd/cli/cli.go
```

#### Run Tests

```shell
govendor test +local
```

#### Saving dependencies

If you make a change and import a new package, run this to vendor the imports.

```shell
govendor add +external
```

#### Compiling Provider files, UI templates, and static assets

Supergiant uses [go-bindata](https://github.com/jteeuwen/go-bindata) to compile
assets directly into the code. You will need to run this command if you're
making changes to the UI _or_ if you're working with Provider code:

```shell
go-bindata -pkg bindata -o bindata/bindata.go config/providers/... ui/assets/... ui/views/...
```

#### Enabling SSL

Our AMI distribution automatically sets up self-signed SSL for Supergiant, but
the default [config/config.json.example](config/config.json.example)
does not enable SSL.

You can see [our AMI boot file](build/sgboot) for an example of how
that is done if you would like to use SSL locally or on your own production
setup.

---

## License

This software is licensed under the Apache License, version 2 ("ALv2"), quoted below.

Copyright 2016 Qbox, Inc., a Delaware corporation. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License"); you may not
use this file except in compliance with the License. You may obtain a copy of
the License at http://www.apache.org/licenses/LICENSE-2.0.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
License for the specific language governing permissions and limitations under
the License.

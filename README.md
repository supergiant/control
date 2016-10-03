SUPERGIANT: Easy container orchestration using Kubernetes
=========================================================

---

<!-- Links -->

[Kubernetes Source URL]: https://github.com/kubernetes/kubernetes
[Supergiant Website URL]: https://supergiant.io/
[Supergiant Docs URL]: https://supergiant.io/docs
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
[Create Admin User Anchor]: #create-an-admin-user
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

* [How to Install Supergiant Container Orchestration Engine on AWS EC2][Tutorial
AWS URL]
* [Deploy a MongoDB Replica Set with Docker and Supergiant][Tutorial MongoDB URL]

---

## Features

* Lets you manage microservices with Docker containers
* Lets you manage multiple users (OAUTH and LDAP coming soon)
* Web dashboard served over HTTPS/SSL by default
* Manages hardware like one, big self-healing resource pool
* Lets you easily scale stateful services and HA volumes on the fly
* Lowers costs by auto-scaling hardware when needed
* Lets you easily adjust RAM and CPU min and max values independently for each
service
* Manages hardware topology organically within configurable constraints


## Resources

* [Supergiant Website][Supergiant Website URL]
* [Docs](https://supergiant.io/docs)
* [Tutorials](https://supergiant.io/tutorials)
* [Slack](https://supergiant.io/slack)
* [Install][Tutorial AWS URL]


## Installation

The current release installs on Amazon Web Services EC2, using a
publicly-available AMI. Other cloud providers and local installation are in
development.

If you want to install Supergiant, follow the [Supergiant Install
Tutorial][Tutorial AWS URL].


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
* [Go][Go URL] version 1.7 or more recent
* [Govendor][Govendor URL] for running tests
* [Supergiant as a Go package][Supergiant Go Package Anchor]
* [A certificate signing request file][Generate CSR Anchor] for serving over
HTTPS

If you are missing any of these, see below to [install or generate
dependencies][Install Dependencies Anchor].

From the supergiant Go package folder (usually
`~/.go/src/github.com/supergiant/supergiant`) perform the following:

#### Initialize Config File

In `config/config.json.example`, there's a sample config file to get you
started. Just duplicate and rename to config.json. Examples below assume the
default configuration for localhost development.

```shell
cp config/config.json.example config/config.json
```

#### Run Supergiant

```shell
go run cmd/server/server.go --config-file config/config.json
```

The default configuration expects HTTP requests on port `8080` and HTTPS
requests on port `8081`. These may be changed in `config/config.json`. Visiting
http://localhost:8080/ will produce an insecure warning page by design.

Access the dashboard at https://localhost:8081/ui/ with [the generated Admin
username and password][Create Admin User Anchor].

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

---

## Installing/Generating Dependencies

#### How to Install Supergiant as a Go Package

```shell
go get github.com/supergiant/supergiant
```

This will install supergiant in your
[Go workspace package directory][Go Remote Packages URL] (usually something like
`~/.go/src/github.com/supergiant/supergiant`). Packages in this directory
function like Git repositories that are aware of their upstream origin. From
here, you can create your own branches and even checkout branches under
development.

#### How to Generate a Certificate Signing Request File

You will need local RSA `.key` and `.pem` files. The default locations are in
the Supergiant Go package tmp folder (usually
`~/.go/src/github.com/supergiant/supergiant`) as `tmp/supergiant.key`,
`tmp/supergiant.pem`. If you wish to customize these locations, you will need
to edit `config/config.json`. The following steps require no config editing.

Set the following env session variables:

```shell
SSL_KEY_FILE=tmp/supergiant.key
SSL_CRT_FILE=tmp/supergiant.pem
```

Generate the `.key` file

```shell
openssl genrsa -out $SSL_KEY_FILE 2048
```

Generate the CSR file. This step will ask you a few questions about the computer
you are using to generate the file.

When you are asked to enter **Common Name (e.g. server FQDN or YOUR name)**,
enter `localhost`. This may be customized in `config/config.json`.

```shell
openssl req -new -x509 -sha256 -key $SSL_KEY_FILE -out $SSL_CRT_FILE -days 3650
```

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

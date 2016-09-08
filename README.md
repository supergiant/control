SUPERGIANT: Easy container orchestration using Kubernetes
=========================================================

---

<!-- Links -->

[Kubernetes]: https://github.com/kubernetes/kubernetes
[Supergiant Site]: https://supergiant.io/
[Tutorial AWS]: https://supergiant.io/blog/how-to-install-supergiant-container-orchestration-engine-on-aws-ec2?utm_source=github
[Tutorial MongoDB]: https://supergiant.io/blog/deploy-a-mongodb-replica-set-with-docker-and-supergiant?urm_source=github
[Community URL]: https://supergiant.io/community
[Contribution Guidelines URL]: http://supergiant.github.io/docs/community/contribution-guidelines.html
[Community and Contributing Anchor]: #community-and-contributing

<!-- Badges -->

[GoReportCard Widget]: https://goreportcard.com/badge/github.com/supergiant/supergiant
[GoReportCard URL]: https://goreportcard.com/report/github.com/supergiant/supergiant
[GoDoc Widget]: https://godoc.org/github.com/supergiant/supergiant?status.svg
[GoDoc URL]: https://godoc.org/github.com/supergiant/supergiant
[Travis Widget]: https://travis-ci.org/supergiant/supergiant.svg?branch=master
[Travis URL]: https://travis-ci.org/supergiant/supergiant
[Release Widget]: https://img.shields.io/github/release/supergiant/supergiant.svg
[Release URL]: https://github.com/supergiant/supergiant/releases/latest
[Swagger API Widget]: http://online.swagger.io/validator?url=http://swagger.supergiant.io/api-docs
[Swagger URL]: http://swagger.supergiant.io/docs/

### <img src="http://supergiant.io/img/logo_dark.svg" width="400">

[![GoReportCard Widget]][GoReportCard URL] [![GoDoc Widget]][GoDoc URL] [![Travis Widget]][Travis URL] [![Release Widget]][Release URL]

---


Supergiant is an open-source container orchestration system that lets developers easily deploy and manage apps as Docker containers using Kubernetes.

We want to make Supergiant the easiest way to run Kubernetes in the cloud.

Quick start...

* [How to Install Supergiant Container Orchestration Engine on AWS EC2][Tutorial AWS]
* [Deploy a MongoDB Replica Set with Docker and Supergiant][Tutorial MongoDB]

---

## Features

* Lets you manage microservices with Docker containers
* Lets you manage multiple users (OAUTH and LDAP coming soon)
* Web dashboard served over HTTPS/SSL by default
* Manages hardware like one, big self-healing resource pool
* Lets you easily scale stateful services and HA volumes on the fly
* Lowers costs by auto-scaling hardware when needed
* Lets you easily adjust RAM and CPU min and max values independently for each service
* Manages hardware topology organically within configurable constraints


## Resources

* [Website](https://supergiant.io/)
* [Docs](https://supergiant.io/docs)
* [Tutorials](https://supergiant.io/tutorials)
* [Slack](https://supergiant.io/slack)
* [Install][Tutorial AWS]


## Installation

The current release installs on Amazon Web Services EC2, using a
publicly-available AMI. Other cloud providers and local installation are in
development.

If you want to install Supergiant, follow the [Supergiant Install Tutorial][Tutorial AWS].


## Top-Level Concepts

[![Swagger API Widget]][Swagger URL]

Supergiant makes it easy to run Dockerized apps as services in the cloud by
abstracting Kubernetes resources. It doesn’t obscure Kubernetes in any way --
in fact you could simply use Supergiant to install Kubernetes.

Supergiant abstracts Kubernetes and cloud provider services into a few
easily-managed resources, namely:

* Apps
* Entrypoints
* Components
* Releases

**Apps** are what groups Components into Kubernetes Namespaces. An App is how to
organize some collection of (micro)services in an environment, such as
"my-app-production.”Organization is flexible and up to the user.

**Entrypoints** allow Components to be reached through a public, internet-facing
address. They are how Supergiant handles external load balancing. Kubernetes
handles internal load balancing among containers brilliantly, so we use
Entrypoints as a more efficient system for external load balancing among Nodes.

**A Component** is child of an App and is synonymous with microservice; in that, a
Component should ideally have one role or responsibility within an App. As a
basic example: within an App named "wordpress-production", there might be two
Components named "mysql" and "wordpress".

**A Release** is a configuration of a Component, released at a certain time.
Releases can be verbose, as they represent an entire topology of a Component,
it’s storage volumes, its min and max allocated resources, etc. By managing
Docker instances as Releases, HA storage volumes can be attached and reattached
without losing statefulness.


## Micro-Roadmap

Currently, the core team is working on the following:

* Add LDAP and OAUTH user authentication
* Add support for additional cloud providers
* Add support for local installations


## Community and Contributing

We are very grateful of any contribution.

All Supergiant projects require familiarization with our Community and our Contribution Guidelines. Please see these links to get started.

* [Community Page][Community URL]
* [Contribution Guidelines][Contribution Guidelines URL]


## Development

If you would like to contribute changes to Supergiant, first see the pages in
the section above, [Community and Contributing][Community and Contributing Anchor].

In order to set up a development environment, do the following:

#### Create Admin User

```shell
godep go run cmd/generate_admin_user/generate_admin_user.go --config-file config/config.json
```

#### Run

```shell
godep go run main.go --config-file config/config.json
```

#### Test

```shell
godep go test -v ./test/...
```


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

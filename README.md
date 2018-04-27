Supergiant: Easy container orchestration using Kubernetes
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

### <img src="http://supergiant.io/img/logo_dark.svg" width="400">

[![GoReportCard Widget]][GoReportCard URL] [![GoDoc Widget]][GoDoc URL] [![Travis Widget]][Travis URL] [![Release Widget]][Release URL] [![Coverage Status Widget]][Coverage Status]

---

Supergiant is an open-source container orchestration system that lets developers easily deploy and manage apps as Docker containers using Kubernetes.

Supergiant aims to automate installation and simplify management of Kubernetes across different cloud accounts by providing an easy-to-use UI that exposes Kubernetes API. The system implements simple concepts that abstract Kubernetes API for pod and services deployment, storage, load-balancing, hardware auto-scaling, and more. Supergiant uses an efficient packing algorithm to enable seamless auto-scaling of Kubernetes clusters to minimize costs and improve resilience of applications. For a more detailed overview of Supergiant top-level concepts, see  [the docs folder](docs/v0/).

## Features

* Fully compatible with native Kubernetes (works with existing setups)
* UI and CLI, both built on top of an API (with importable [Go client lib](pkg/client))
* Launch and manage multiple Kubes across multiple cloud providers from the UI
* Works with multiple cloud providers (AWS, DigitalOcean, OpenStack, Packet.net, and
  _actively_ adding more, in addition to on-premise hardware support)
* Deploy / Update / Restart containers with a few clicks
* Get instant access to apps from Helm Repositories
* Filterable container metrics views (RAM / CPU timeseries graphs)
* Get a comprehensive view of cluster resources with built-in monitoring and logging
* Automatic server management (manual addition of new nodes, background server autoscaling, up/down depending on container resource needs)
* Role-based Users, Session-based login, self-signed SSL, and API tokens for
  security (OAuth and LDAP support soon coming)


## Installation

1. Download the Supergiant server for your system (Windows, Mac, and Linux) and processor architecture from our  [releases](https://github.com/supergiant/supergiant/releases) page. For example, for Linux: 

```sh
curl https://github.com/supergiant/supergiant/releases/download/v0.15.6/supergiant-server-linux-amd64 -L -o /usr/bin/supergiant
```

2. Make sure to make the downloaded binary executable:

```sh
sudo chmod +x /usr/bin/supergiant
```

3. Download  [the example config file](https://github.com/supergiant/supergiant/blob/master/config/config.json.example) and customize it: 

```sh
curl https://raw.githubusercontent.com/supergiant/supergiant/master/config/config.json.example --create-dirs -o /etc/supergiant/config.json
```

In the configuration file, specify your paths to log and database locations and create corresponding directories for them.

```json
{
 ...
 "sqlite_file": "/var/lib/supergiant/development.db",
 ...
 "log_file": "/var/log/supergiant/development.log",
 ...
}
```

4. Run the binary with a config file and save the user/password for the admin user generated the first time Supergiant runs.

```json
<supergiant-server-binary> --config-file /etc/supergiant/config.json
```

5. Access Supergiant on default 8080 port on localhost. 

#### Installing on AWS

![](https://supergiant.io/uploads/blog/sg_aws_step_1@2x.jpg)

If you want to easily install Supergiant on Amazon Web Services EC2 with Supergiant Amazon Machine Image (AMI), follow the [Supergiant AWS Install Tutorial](https://supergiant.io/blog/how-to-install-supergiant-container-orchestration-engine-on-aws-ec2?utm_source=github).

## Usage

#### Deploying Kubernetes Cluster

 Supergiant allows deploying Kubernetes clusters (Kubes) via the easy-to-use interface with the minimal configuration required. The system manages Kubernetes installation and configuration under the hood enabling various master and nodes services and tools. 

![](http://res.cloudinary.com/doj9feked/image/upload/v1523603021/kube-deploy_ibkfcd.gif)

#### Deploying Apps from Helm Repositories

Supergiant provides access to Kubernetes curated helm charts with a broad choice of configured containers installable in one click. 

![](http://res.cloudinary.com/doj9feked/image/upload/v1523603035/app-deploy_nnsy5t.gif)

#### Tracking Cluster Resources

Supergiant UI gives a comprehensive view of computer resources used by the cluster, running applications, services, and attached storages to get insights and simplify cluster administration.

![](http://res.cloudinary.com/doj9feked/image/upload/v1523605692/cluster-resources_ylbset.gif)



## Micro-Roadmap

Currently, the core team is working on the following:

* Add LDAP and OAuth user authentication
* Add support for new cloud providers
* Add support for local installations


## Resources

- [Supergiant Website][Supergiant Website URL]
- [Top-level concepts](docs/v0/)
- [Tutorials](https://supergiant.io/tutorials)
- [Slack Support Channel](https://supergiant.io/slack)
- [Install on AWS][Tutorial AWS URL]

## Community and Contributing

We are grateful for any contribution to the Supergiant project be it in a form of a new GitHub issue, a GitHub feature Pull Request, social media engagement etc. Contributing to Supergiant projects requires familiarization with Community and our Contribution Guidelines. Please see these links to get started.

* [Community Page][Supergiant Community URL]
* [Contribution Guidelines][Supergiant Contribution Guidelines URL]


## Development

#### Use Docker in development

    docker-compose build server
    docker-compose run --rm --service-ports server

#### Native go on your host

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

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under
the License.

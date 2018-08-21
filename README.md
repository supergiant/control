Supergiant: Kubernetes Orchestration
===

<!-- Badge Links -->
[GoReportCard Widget]: https://goreportcard.com/badge/github.com/supergiant/supergiant
[GoReportCard URL]: https://goreportcard.com/report/github.com/supergiant/supergiant
[GoDoc Widget]: https://godoc.org/github.com/supergiant/supergiant?status.svg
[GoDoc URL]: https://godoc.org/github.com/supergiant/supergiant
[Travis Widget]: https://travis-ci.org/supergiant/supergiant.svg?branch=master
[Travis URL]: https://travis-ci.org/supergiant/supergiant
[Release Widget]: https://img.shields.io/github/release/supergiant/supergiant.svg
[Release URL]: https://github.com/supergiant/supergiant/releases/latest
[Coverage Status]: https://coveralls.io/github/supergiant/supergiant?branch=master
[Coverage Status Widget]: https://coveralls.io/repos/github/supergiant/supergiant/badge.svg?branch=master

# <img src="http://supergiant.io/img/logo_dark.svg" width="400">

<!-- Badges -->
[![Release Widget]][Release URL] [![GoDoc Widget]][GoDoc URL] [![Travis Widget]][Travis URL] [![Coverage Status Widget]][Coverage Status] [![GoReportCard Widget]][GoReportCard URL]

---

Supergiant empowers developers and administrators through its simplified deployment and management of Kubernetes, in addition to easing the configuration and deployment of Helm charts, taking advantage of Kubernetes' power, flexibility, and abstraction.

Supergiant facilitates clusters on multiple cloud providers, striving for truly agnostic and impartial infrastructure--and it does this with an autoscaling system that cares deeply about efficiency. It asserts through downscaling and resource compaction that unutilized infrastructure shouldn't be paid for (and, therefore, shouldn't be running).

Supergiant implements simple practices that abstract load-balancing, application deployment, basic monitoring, node deployment or destruction, and more, on a highly usable UI. Its efficient compaction algorithm enables seamless auto-scaling of Kubernetes clusters, minimizing costs while maintaining the resiliency of applications. To dive into top-level concepts, see [the documentation](https://supergiant.readme.io/v1.1.0/docs).

# Features

* Fully compatible with native Kubernetes versions 1.5.7, 1.6.7, 1.7.7, and 1.8.7
* Easy management and deployment of multiple kubes in various configurations
* AWS, DigitalOcean, OpenStack, Packet, GCE, and on-premise kube deployment
* Easy creation of Helm releases, Pods, Services, LoadBalancers, etc.
* Automatic, resource-based node scaling
* Compatibility with multiple hardware architectures
* Role-based Users, Session-based logins, self-signed SSLs, and API tokens
* A clean UI and CLI, both built on top of an API (with importable [Go client lib](pkg/client))

# Resources

- [Supergiant Website](https://supergiant.io/)
- [Top-level concepts](https://supergiant.readme.io/v1.1.0/docs/architecture-purpose)
- [Tutorials](https://supergiant.io/tutorials)
- [Slack Support Channel](https://supergiant.io/slack)
- [Installation](https://supergiant.readme.io/v1.1.0/docs/installation)
~ [UI Usage](https://supergiant.readme.io/v1.1.0/docs/using-the-ui)
~ [API Usage](https://supergiant.readme.io/v1.1.0/docs/using-the-api)

# Community and Contributing

We are grateful for any contribution to the Supergiant project, be it in the form of a GitHub issue, a Pull Request, or any social media engagement. Contributing to Supergiant requires familiarization with the Community and Contribution Guidelines. Please see these links to get started:

* [Community Page](https://supergiant.io/community)
* [Contribution Guidelines](http://supergiant.github.io/docs/community/contribution-guidelines.html)

## Development

### Setup

Note: _These instructions are for setting up a local or custom environment and assume a "beginner" level of experience. For executable binaries, see the [Releases Page](https://github.com/supergiant/supergiant/releases). Before submitting a Pull Request, please see [Community and Contributing](#community-and-contributing) above._

#### 1. Install Dependencies:

* [Git](https://git-scm.com/)
* [Go](https://golang.org/) - version 1.7+ 
* [govendor](https://github.com/kardianos/govendor)
* [go-bindata](https://github.com/go-bindata/go-bindata)
* [npm & Node.js](https://www.npmjs.com/get-npm)
* [Angular](https://cli.angular.io/) - version 5.0+

Note: _If a new package is imported and then used in Supergiant code, make sure to vendor the imports (the govendor binary is located in `$GOPATH/bin/govendor`):_

```shell
govendor add +external
```

Note: _New to any of these technologies? Here are a few basic, free resources:_
* An [introduction to Git](https://www.youtube.com/watch?v=xuB1Id2Wxak)
* A [Golang tutorial playlist](https://www.youtube.com/watch?v=G3PvTWRIhZA&list=PLQVvvaa0QuDeF3hP0wQoSxpkqgRcgxMqX)
* A [guide to using npm](https://www.youtube.com/watch?v=jHDhaSSKmB0)
* A [guide to Angular 5](https://www.youtube.com/watch?v=AaNZBrP26LQ)

#### 2. Fork the Repo

Please make sure to fork the Supergiant project before getting started. Direct Pull Requests will not be accepted into the project.

To make the commands in subsequent steps copy-pastable, feel free to do the folowing:

```shell
GITHUB_USERNAME="github_username"
```

#### 3. Clone the Fork

Clone the repo into a local Go project directory, and add an upstream remote to keep the fork synced:

```shell
git clone https://github.com/$GITHUB_USERNAME/supergiant.git $GOPATH/src/github.com/supergiant/supergiant
git remote add upstream https://github.com/supergiant/supergiant.git
```

Remember to checkout a branch to work on, or create a new branch.

#### 4. Setup a Config File

Note: _From now on, all specified directories have a root of `$GOPATH/src/github.com/supergiant/`._

Copy the [example configuration](https://github.com/supergiant/supergiant/blob/v1.0.0/config/config.json.example) from within `./supergiant/`, or [create a custom one](https://supergiant.readme.io/v1.1.0/docs/configuration) as desired:

```shell
cp config/config.json.example config/config.json
```

#### 5. Compile UI Assets

Run this from within `./supergiant/`. It was once used to run the old UI, and, although the old UI is no longer in use, it is necessary to run the backend (the binary for go-bindata is found in `$GOPATH/bin/`):

```shell
go-bindata -pkg bindata -o bindata/bindata.go config/providers/... ui/assets/... ui/views/...
```

### Operation

#### 1. Running the Backend

From within `./supergiant/`, run:

```shell
go run cmd/server/server.go --config-file config/config.json
```

The server will output as seen below:

```shell
INFO[0000] No Admin detected, creating new and printing credentials:

  ( ͡° ͜ʖ ͡°)  USERNAME: admin  PASSWORD: 0GYB4rIU8TaokNtJ
```

Note: _The initial account login credentials will appear like this the first time that the server runs, and will not appear again, though anyone with access to these logs can see them. The credentials can be changed._

#### 2. Running the UI

Currently, Supergiant uses a UI developed with Angular 5 and higher. The Angular project directory is found in `./supergiant/cmd/ui/assets`. The old UI is accessible on port `8080` when the Supergiant Server is running.

##### 2.a. Install Node Modules

In `./supergiant/cmd/ui/assets`, run:

```shell
npm install
```

Note: _If the UI fails to initialize, package and dependency problems are often the cause. If there is a problem with the project's current `package.json` or `package-lock.json`, please open a GitHub issue._

##### 2.b. Serve the UI

Within `./supergiant/cmd/ui/assets/`, run:

```shell
ng serve
```

The UI will be accessible on port 4200 by default. The server must be running to interact properly with the UI.

### Building

#### 1. Building the Server

Note: _The build processes for both the server binary and Docker image are currently under review._

#### 2. Building the UI

Run `ng build` to build the project. The build artifacts will be stored in the `dist/` directory. Use the `-prod` flag for a production build.

Note: _The build processes for both the UI binary and Docker image are currently under review._

### Testing

#### 1. Testing the Server

Currently, tests are split up into two groups: unit tests and integration tests. Unit tests are in `./supergiant/pkg/`, and integration tests are in `./supergiant/test`.

##### 1.a. Unit Testing

To run unit tests, simply run `go test` on all of the Supergiant-specific packages within `./supergiant/`:

```shell
go test -v ./pkg/...
```

##### 1.b. Integration Testing

Note: _Before running integration tests, make sure to set the required environment variables noted in the code for each provider in `./supergiant/test/providers/`_

For integration tests, Golang build tags are used. This helps with separation of testing concerns. The tag used is `// +build integration`. Like unit tests, it is easy to run the tests from `./supergiant/`:

```shell
go test -v -tags=integration ./test/...
```

#### 2. Testing the UI

##### 1.a. Unit Testing

Run `ng test` to execute the unit tests with [Karma](https://karma-runner.github.io).

##### 1.b. E2E Testing

Run `ng e2e` to execute the end-to-end tests with [Protractor](http://www.protractortest.org/).
Before running the tests, make sure to serve the app via `ng serve`.

<!-- ### Enabling SSL TODO: Figure out if/how we do this. -->

## License

Copyright 2016 Qbox, Inc., a Delaware corporation. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License"); you may not
use this file except in compliance with the License. You may obtain a copy of
the License at http://www.apache.org/licenses/LICENSE-2.0.

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

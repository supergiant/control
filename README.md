[![GoReportCard Widget]][GoReportCard] [![GoDoc Widget]][GoDoc] [![Travis Widget]][Travis] [![Coverage Status Widget]][Coverage Status]

[GoReportCard Widget]: https://goreportcard.com/badge/github.com/supergiant/supergiant
[GoReportCard]: https://goreportcard.com/report/github.com/supergiant/supergiant
[GoDoc]: https://godoc.org/github.com/supergiant/supergiant
[GoDoc Widget]: https://godoc.org/github.com/supergiant/supergiant?status.svg
[Travis]: https://travis-ci.org/supergiant/supergiant
[Travis Widget]: https://travis-ci.org/supergiant/supergiant.svg?branch=master
[Coverage Status]: https://coveralls.io/github/supergiant/supergiant?branch=master
[Coverage Status Widget]: https://coveralls.io/repos/github/supergiant/supergiant/badge.svg?branch=master

# Supergiant

Supergiant is API-based stateful container orchestration disguised as a
developer-friendly application platform. It is based on
[Kubernetes](https://github.com/kubernetes/kubernetes).

[supergiant.io](https://supergiant.io)

### Concepts

See the [docs](docs/v0/apps.md).

### Development

```shell
# Run etcd in background (or in another tab/pane)
etcd &

go run main.go \
--etcd-hosts http://localhost:2379 \
--aws-region us-east-1 \
--aws-az us-east-1c \
--aws-sg-id <security_group_id> \
--aws-subnet-id <subnet_id> \
--k8s-host <kubernetes_master_host> \
--k8s-user <kubernetes_api_user> \
--k8s-pass <kubernetes_api_pass> \
--aws-access-key <AWS access key> \
--aws-secret-key <AWS secret key> \
--k8s-insecure-https \
--enable-capacity-service \
--log-level=debug
```

### Tests

*Better coverage is on the way, we promise.*

```shell
godep go test ./...
```

### License

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

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
[Supergiant Swagger Docs URL]: http://swagger.supergiant.io/docs/
[Tutorial AWS URL]: https://supergiant.io/blog/how-to-install-supergiant-container-orchestration-engine-on-aws-ec2?utm_source=github
[Tutorial MongoDB URL]: https://supergiant.io/blog/deploy-a-mongodb-replica-set-with-docker-and-supergiant?urm_source=github
[Community and Contributing Anchor]: #community-and-contributing
[Swagger URL]: http://swagger.io/
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
[Swagger API Widget]: http://online.swagger.io/validator?url=http://swagger.supergiant.io/api-docs
[Swagger URL]: http://swagger.supergiant.io/docs/

### <img src="http://supergiant.io/img/logo_dark.svg" width="400">

[![GoReportCard Widget]][GoReportCard URL] [![GoDoc Widget]][GoDoc URL] [![Travis Widget]][Travis URL] [![Release Widget]][Release URL]

---

## Development

If you would like to contribute changes to Supergiant, first see the pages in the Community and Contributing section of the [README.md](https://github.com/supergiant/supergiant/blob/master/README.md#community-and-contributing).

To compile the UI assets you will need:

  - [Sass]: http://sass-lang.com/

You can compile the Sass assets with the following terminal command.

```shell
sass _bootstrap.scss main.css --style compressed
```


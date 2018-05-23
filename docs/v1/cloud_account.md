# Cloud Account

A Cloud Account is a simple object that holds API access credentials (e.g., API Token, security keys) for any supported Provider (e.g., GCE, AWS, Packet). It is the parent object of [Kubes](kube.md).  Credentials are validated by an API call to the respective Provider.

### DigitalOcean

![](../img/digital-ocean-credentials.png)

### Amazon AWS

![](../img/aws-cloud-acc-credentials.png)

### OpenStack

![](../img/openstack-cloud-credentials.png)



See **[Linking Cloud Accounts](https://github.com/supergiant/supergiant/wiki/Linking-Cloud-Accounts)** guide for a detailed information about linking cloud accounts in Supergiant 1.0.0
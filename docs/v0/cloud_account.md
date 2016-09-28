# Cloud Account

A Cloud Account is a simple object that holds API access credentials for any
supported Provider.

It is the parent object of [Kubes](kube.md). Credentials are validated
on create by an API call to the respective Provider.

### Examples

#### AWS

```json
{
  "name": "my-aws-account",
  "provider": "aws",
  "credentials": {
    "access_key": "<your_access_key>",
    "secret_key": "<your_secret_key>"
  }
}
```

#### DigitalOcean

```json
{
  "name": "my-do-account",
  "provider": "digitalocean",
  "credentials": {
    "token": "<your_api_token>"
  }
}
```

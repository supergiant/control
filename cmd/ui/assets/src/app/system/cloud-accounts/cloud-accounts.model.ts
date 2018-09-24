// This is the only model that is pre-populated with providers given that the possibilities are known in advance
// TODO: SO much duplicated data
export class CloudAccountModel {
  aws = {
    'name' : 'AWS - Amazon Web Services',
    'model': {
      'provider': 'aws',
      'credentials': {
        'access_key': '',
        'secret_key': ''
      }
    },
    'schema': {
      'properties': {
        'provider': {
          'default': 'aws',
          'description': 'AWS - Amazon Web Services',
          'type': 'string',
          'widget': 'hidden'
        },
        'credentials': {
          'type': 'object',
          'properties': {
            'access_key': {
              'type': 'string',
              'description': 'IAM user access key ID'
            },
            'secret_key': {
              'type': 'string',
              'description': 'IAM user secret access key'
            }
          }
        }
      }
    },
    'layout': [
      { "key": "credentials.access_key", "placeholder": "Access Key" },
      { "key": "credentials.secret_key", "placeholder": "Secret Key" },
      { "type": "submit", "title": "CREATE" }
    ]
  };

  digitalocean = {
    'name' : 'Digital Ocean',
    'model': {
      'provider': 'digitalocean',
      'credentials': {
        'accessToken': '',
        'publicKey': ''
      }
    },
    'schema': {
      'properties': {
        'provider': {
          'default': 'digitalocean',
          'description': 'Digital Ocean',
          'type': 'string',
          'widget': 'hidden'
        },
        'credentials': {
          'type': 'object',
          'properties': {
            'accessToken': {
              'type': 'string',
              'description': 'Access Token for your DO account'
            },
            'publicKey': {
              'type': 'string',
              'description': 'Your personal private key'
            }
          }
        }
      }
    },
    'layout': [
      { "key": "credentials.accessToken", "placeholder": "Access Token" },
      { "type": "textarea", "key": "credentials.publicKey", "placeholder": "Public Key" },
      { "type": "submit", "title": "CREATE" }
    ]
  };

  gce = {
    'name' : 'GCE - Google Compute Engine',
    'model': {
      'provider': 'gce',
      'credentials': {
        'service_account_key': 'Paste your Service Account Key here...'
      }
    }
  };

  packet = {
    'name' : 'Packet.net',
    'model': {
      'provider': 'packet',
      'credentials': {
        'api_token': ''
      }
    },
    'schema': {
      'properties': {
        'provider': {
          'default': 'packet',
          'description': 'Packet.net',
          'type': 'string',
          'widget': 'hidden'
        },
        'credentials': {
          'type': 'object',
          'properties': {
            'api_token': {
              'type': 'string',
              'description': 'API Token for your Packet account'
            }
          }
        }
      }
    },
    'layout': [
      { "key": "credentials.api_token", "placeholder": "API Token" },
      { "type": "submit", "title": "CREATE" }
    ]
  };

  public providers = [
    {
      display: "AWS - Amazon Web Services",
      name: "aws",
      data: this.aws
    },
    {
      display: "Digital Ocean",
      name: "digitalocean",
      data: this.digitalocean
    },
    {
      display: "GCE - Google Compute Engine",
      name: "gce",
      data: this.gce
    },
    {
      display: "Packet.net",
      name: "packet",
      data: this.packet
    }
  ]
}

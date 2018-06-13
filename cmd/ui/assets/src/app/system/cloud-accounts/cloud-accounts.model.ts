// This is the only model that is pre-populated with providers given that the possibilities are known in advance
export class CloudAccountModel {
  aws = {
    'name' : 'AWS - Amazon Web Services',
    'model': {
      'credentials': {
        'access_key': '',
        'secret_key': ''
      },
      'name': '',
      'provider': 'aws'
    },
    'schema': {
      'properties': {
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
        },
        'name': {
          'type': 'string',
          'description': 'Choose a name for this cloud account'
        },

        'provider': {
          'default': 'aws',
          'description': 'AWS - Amazon Web Services',
          'type': 'string',
          'widget': 'hidden'
        }
      }
    },
    'layout': [
      { 'key': 'name' },
      {
        'key': 'credentials',
        'type': 'div',
        'htmlClass': 'creds-wrapper',
        'items': [
          { 'key': 'credentials.access_key' },
          { 'key': 'credentials.secret_key' }
        ]
      },
      {
        'type': 'submit',
        'title': 'Create'
      }
    ]
  };

  digitalocean = {
    'name' : 'Digital Ocean',
    'model': {
      'credentials': {
        'token': ''
      },
      'name': '',
      'provider': 'digitalocean'
    },
    'schema': {
      'properties': {
        'credentials': {
          'type': 'object',
          'properties': {
            'token': {
              'description': 'API Token',
              'type': 'string'
            }
          }
        },
        'name': {
          'description': 'Choose a name for this cloud account',
          'type': 'string'
        },
        'provider': {
          'default': 'digitalocean',
          'description': 'Digital Ocean',
          'type': 'string',
          'widget': 'hidden'
        }
      }
    },
    'layout': [
      { 'key': 'name' },
      {
        'key': 'credentials',
        'type': 'div',
        'htmlClass': 'creds-wrapper',
        'items': [
          { 'key': 'credentials.token' }
        ]
      },
      {
        'type': 'submit',
        'title': 'Create'
      }
    ]
  };

  gce = {
    'name' : 'GCE - Google Compute Engine',
    'model': {
      'credentials': {
        'service_account_key': ''
      },
      'name': '',
      'provider': 'gce'
    },
    'schema': {
      'properties': {
        'credentials': {
          'type': 'object',
          'properties': {
            'service_account_key': {
              'description': 'Paste in your service account key here.',
              'type': 'string'
            }
          }
        },
        'name': {
          'description': 'Choose a name for this cloud account',
          'type': 'string'
        },
        'provider': {
          'default': 'gce',
          'description': 'GCE - Google Compute Engine',
          'type': 'string',
          'widget': 'hidden'
        }
      }
    },
    'layout': [
      { 'key': 'name' },
      {
        'key': 'credentials',
        'type': 'div',
        'htmlClass': 'creds-wrapper',
        'items': [
          {
            'key': 'credentials.service_account_key',
            'type': 'textarea',
          }
        ]
      },
      {
        'type': 'submit',
        'title': 'Create'
      }
    ]
  };

  openstack = {
    'name' : 'OpenStack',
    'model': {
      'credentials': {
        'domain_id': '',
        'domain_name': '',
        'identity_endpoint': '',
        'password': '',
        'tenant_id': '',
        'username': ''
      },
      'name': '',
      'provider': 'openstack'
    },
    'schema': {
      'properties': {
        'credentials': {
          'type': 'object',
          'properties': {
            'domain_id': {
              'description': 'Domain ID',
              'type': 'string'
            },
            'domain_name': {
              'description': 'Domain Name',
              'type': 'string'
            },
            'identity_endpoint': {
              'description': 'Identity Endpoint',
              'type': 'string'
            },
            'password': {
              'description': 'Password',
              'type': 'string'
            },
            'tenant_id': {
              'description': 'Tenant ID',
              'type': 'string'
            },
            'username': {
              'description': 'User Name',
              'type': 'string'
            }
          }
        },
        'name': {
          'description': 'Choose a name for this cloud account',
          'type': 'string'
        },
        'provider': {
          'default': 'openstack',
          'description': 'OpenStack',
          'type': 'string',
          'widget': 'hidden'
        }
      }
    },
    'layout': [
      { 'key': 'name' },
      {
        'key': 'credentials',
        'type': 'div',
        'htmlClass': 'creds-wrapper',
        'items': [
          { 'key': 'credentials.domain_id' },
          { 'key': 'credentials.domain_name' },
          { 'key': 'credentials.identity_endpoint' },
          { 'key': 'credentials.password' },
          { 'key': 'credentials.tenant_id' },
          { 'key': 'credentials.username' }
        ]
      },
      {
        'type': 'submit',
        'title': 'Create'
      }
    ]
  };

  packet = {
    'name' : 'Packet.net',
    'model': {
      'credentials': {
        'token': ''
      },
      'name': '',
      'provider': 'packet'
    },
    'schema': {
      'properties': {
        'credentials': {
          'type': 'object',
          'properties': {
            'token': {
              'description': 'API Token',
              'type': 'string'
            }
          }
        },
        'name': {
          'description': 'Choose a name for this cloud account',
          'type': 'string'
        },
        'provider': {
          'default': 'packet',
          'description': 'Packet.net',
          'type': 'string',
          'widget': 'hidden'
        }
      }
    },
    'layout': [
      { 'key': 'name' },
      {
        'key': 'credentials',
        'type': 'div',
        'htmlClass': 'creds-wrapper',
        'items': [
          { 'key': 'credentials.token' }
        ]
      },
      {
        'type': 'submit',
        'title': 'Create'
      }
    ]
  };

  public providers = [
    {
      name: 'AWS - Amazon Web Services',
      data: this.aws,
    },
    {
      name: 'Digital Ocean',
      data: this.digitalocean,
    },
    {
      name: 'GCE - Google Compute Engine',
      data: this.gce,
    },
    {
      name: 'OpenStack',
      data: this.openstack,
    },
    {
      name: 'Packet.net',
      data: this.packet,
    },
  ];
}

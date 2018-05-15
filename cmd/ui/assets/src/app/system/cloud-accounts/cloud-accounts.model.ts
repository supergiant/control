// This is the only model that is pre-populated with providers given that the possibilities are known in advance
export class CloudAccountModel {
  aws = {
    'name': 'AWS - Amazon Web Services',
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
          'properties': {
            'access_key': {
              'description': 'Access Key',
              'type': 'string'
            },
            'secret_key': {
              'description': 'Secret Access Key',
              'type': 'string'
            }
          },
          'type': 'object'
        },
        'name': {
          'description': 'Provider Name',
          'type': 'string'
        },
        'provider': {
          'default': 'aws',
          'description': 'AWS - Amazon Web Services',
          'type': 'string',
          'widget': 'hidden'
        }
      }
    }
  };
  digitalocean = {
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
          'properties': {
            'token': {
              'description': 'API Token',
              'type': 'string'
            }
          },
          'type': 'object'
        },
        'name': {
          'description': 'Provider Name',
          'type': 'string'
        },
        'provider': {
          'default': 'digitalocean',
          'description': 'Digital Ocean',
          'type': 'string',
          'widget': 'hidden'
        }
      }
    }
  };

  gce = {
    'model': {
      'credentials': {
        'auth_provider_x509_cert_url': '',
        'auth_uri': '',
        'client_email': '',
        'client_id': '',
        'client_x509_cert_url': '',
        'private_key': '',
        'private_key_id': '',
        'project_id': '',
        'token_uri': '',
        'type': ''
      },
      'name': '',
      'provider': 'gce'
    },
    'schema': {
      'properties': {
        'credentials': {
          'properties': {
            'auth_provider_x509_cert_url': {
              'description': 'auth_provider_x509_cert_url',
              'type': 'string'
            },
            'auth_uri': {
              'description': 'Auth URI',
              'type': 'string'
            },
            'client_email': {
              'description': 'Client Email',
              'type': 'string'
            },
            'client_id': {
              'description': 'Client ID',
              'type': 'string'
            },
            'client_x509_cert_url': {
              'description': 'client_x509_cert_url',
              'type': 'string'
            },
            'private_key': {
              'description': 'Private Key',
              'type': 'string'
            },
            'private_key_id': {
              'description': 'Private Key ID',
              'type': 'string'
            },
            'project_id': {
              'description': 'Project ID',
              'type': 'string'
            },
            'token_uri': {
              'description': 'Token URI',
              'type': 'string'
            },
            'type': {
              'description': 'Type',
              'type': 'string'
            }
          },
          'type': 'object'
        },
        'name': {
          'description': 'Provider Name',
          'type': 'string'
        },
        'provider': {
          'default': 'gce',
          'description': 'GCE - Google Compute Engine',
          'type': 'string',
          'widget': 'hidden'
        }
      }
    }
  };
  openstack = {
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
          },
          'type': 'object'
        },
        'name': {
          'description': 'Provider Name',
          'type': 'string'
        },
        'provider': {
          'default': 'openstack',
          'description': 'OpenStack',
          'type': 'string',
          'widget': 'hidden'
        }
      }
    }
  };
  packet = {
    'model': {
      'credentials': {
        'api_token': ''
      },
      'name': '',
      'provider': 'packet'
    },
    'schema': {
      'properties': {
        'credentials': {
          'properties': {
            'api_token': {
              'description': 'API Token',
              'type': 'string'
            }
          },
          'type': 'object'
        },
        'name': {
          'description': 'Provider Name',
          'type': 'string'
        },
        'provider': {
          'default': 'packet',
          'description': 'Packet.net',
          'type': 'string',
          'widget': 'hidden'
        }
      }
    }
  };
  public providers = {
    'AWS - Amazon Web Services': this.aws,
    'Digital Ocean': this.digitalocean,
    'GCE - Google Compute Engine': this.gce,
    'OpenStack': this.openstack,
    'Packet.net': this.packet
  };
}

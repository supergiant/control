export class ClusterDigitalOceanModel {
  digitalocean = {
    'data': {
      'cloud_account_name': '',
      'digitalocean_config': {
        'region': 'nyc1',
        'kube_master_count': 1,
        'ssh_key_fingerprint': []
      },
      'master_node_size': '1gb',
      'name': '',
      'node_sizes': [
        '1gb',
        '2gb',
        '4gb',
        '8gb',
        '16gb',
        '32gb',
        '48gb',
        '64gb'
      ]
    },
    'schema': {
      'properties': {
        'cloud_account_name': {
          'description': 'Cloud Account Name',
          'type': 'string'
        },
        'digitalocean_config': {
          'properties': {
            'region': {
              'default': 'nyc1',
              'description': 'Region',
              'type': 'string'
            },
           'ssh_key_fingerprint': {
             'description': 'SSH Key Fingerprint',
             'id': '/properties/ssh_key_fingerprint',
             'items': {
               'id': '/properties/ssh_key_fingerprint/items',
               'type': 'string'
             },
             'type': 'array'
           }
          },
          'type': 'object'
        },
        'master_node_size': {
          'default': '1gb',
          'description': 'Master Node Size',
          'type': 'string'
        },
        'name': {
          'description': 'Name',
          'type': 'string',
          'pattern': '^[a-z]([-a-z0-9]*[a-z0-9])?$',
          'maxLength': 12
        },
        'kube_master_count': {
          'description': 'Kube Master Count',
          'type': 'number',
          'widget': 'number',
        },
        'node_sizes': {
          'description': 'Node Sizes',
          'id': '/properties/node_sizes',
          'items': {
            'id': '/properties/node_sizes/items',
            'type': 'string'
          },
          'type': 'array'
        }
      }
    }
  };
}

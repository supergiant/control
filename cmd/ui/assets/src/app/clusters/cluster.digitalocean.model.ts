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
        'name': {
          'description': 'The desired name of the kube. Max length of 12 characters.',
          'type': 'string',
          'pattern': '^[a-z]([-a-z0-9]*[a-z0-9])?$',
          'maxLength': 12
        },
        'kubernetes_version': {
          'default': '1.5.7',
          'description': 'The Version of Kubernetes to be deployed.',
          'type': 'string',
          'enum': ['1.5.7', '1.6.7', '1.7.7', '1.8.7'] // TODO: <-- Should be dynamically populated.
        },
        'cloud_account_name': {
          'description': 'The Supergiant cloud account you created for use with Packet.',
          'type': 'string'
        },
        'digitalocean_config': {
          'properties': {
            'region': {
              'default': 'nyc1',
              'description': 'The Digital Ocean region the kube will be created in.',
              'type': 'string'
            },
            'ssh_key_fingerprint': {
              'description': 'The fingerprint of the public key that you uploaded to your OpenStack account.',
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
          'description': 'The size of the server the master will live on.',
          'type': 'string'
        },
        'kube_master_count': {
          'description': 'The number of masters desired--for High Availability.',
          'type': 'number',
          'widget': 'number',
        },
        'node_sizes': {
          'description': 'The sizes you want to be available to Supergiant when scaling.',
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

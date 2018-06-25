export class ClusterDigitalOceanModel {
  digitalocean = {
    'data': {
      'cloud_account_name': '',
      'digitalocean_config': {
        'region': 'nyc1',
        'ssh_key_fingerprint': []
      },
      'kubernetes_version': '1.8.7',
      'kube_master_count': 1,
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
          'type': 'string',
          'description': 'The desired name of the kube. Max length of 12 characters.',
          'pattern': '^[a-z]([-a-z0-9]*[a-z0-9])?$',
          'maxLength': 12
        },
        'kubernetes_version': {
          'type': 'string',
          'description': 'The Version of Kubernetes to be deployed.',
          'default': '1.5.7',
          'enum': ['1.5.7', '1.6.7', '1.7.7', '1.8.7'] // TODO: <-- Should be dynamically populated.
        },
        'cloud_account_name': {
          'type': 'string',
          'description': 'The Supergiant cloud account you created for use with Digital Ocean.'
        },
        'digitalocean_config': {
          'type': 'object',
          'properties': {
            'region': {
              'type': 'string',
              'description': 'The Digital Ocean region the kube will be created in.',
              'default': 'nyc1',
              'enum': [
                'ams2',
                'ams3',
                'blr1',
                'fra1',
                'lon1',
                'nyc1',
                'nyc2',
                'nyc3',
                'sfo1',
                'sfo2',
                'sgp1',
                'tor1'
              ]
            },
            'ssh_key_fingerprint': {
              'type': 'array',
              'description': 'The fingerprint of the public key that you uploaded to your Digital Ocean account.',
              'id': '/properties/ssh_key_fingerprint',
              'items': {
                'type': 'string',
                'id': '/properties/ssh_key_fingerprint/items'
              }
            }
          },
          'required': [ 'region' ]
        },
        'master_node_size': {
          'type': 'string',
          'description': 'The size of the server the master will live on.',
          'default': '1gb'
        },
        'kube_master_count': {
          'type': 'number',
          'description': 'The number of masters desired - for High Availability.',
          'widget': 'number',
          'minimum': 1
        },
        'node_sizes': {
          'type': 'array',
          'description': 'The sizes you want to be available to Supergiant when scaling.',
          'id': '/properties/node_sizes',
          'items': {
            'type': 'string',
            'id': '/properties/node_sizes/items'
          }
        }
      },
      'required': [ 'name' ]
    },
    'layout': [
      {
        'type': 'section',
        'title': 'Cluster Details',
        'items': [
          { 'key': 'name' },
          { 'key': 'kubernetes_version' },
          // { 'key': 'cloud_account_name', 'readonly': true }
        ]
      },
      {
        'type': 'section',
        'title': 'Digital Ocean Config',
        'items': [
          { 'key': 'digitalocean_config.region' },
          {
            'key': 'digitalocean_config.ssh_key_fingerprint',
            'type': 'array',
            'title': 'SSH Key Fingerprint',
            'items': [
              { 'key': 'digitalocean_config.ssh_key_fingerprint[]' }
            ]
          } // TODO: can you really add multiple fingerprints?
        ]
      },
      {
        'type': 'section',
        'items': [
          { 'key': 'master_node_size' },
          { 'key': 'kube_master_count' },
        ]
      },
      {
        'key': 'node_sizes',
        'type': 'array',
        'items': [
          { 'key': 'node_sizes[]' },
        ]
      },
      {
        'type': 'submit',
        'title': 'Create'
      }
    ]
  };
}

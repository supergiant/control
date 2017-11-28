export class ClusterGCEModel {
  gce = {
    'data': {
      'cloud_account_name': '',
      'gce_config': {
        'ssh_pub_key': '',
        'zone': 'us-east1-b'
      },
      'master_node_size': 'n1-standard-1',
      'name': '',
      'kube_master_count': 1,
      'node_sizes': [
        'n1-standard-1',
        'n1-standard-2',
        'n1-standard-4',
        'n1-standard-8'
      ]
    },
    'schema': {
      'properties': {
        'cloud_account_name': {
          'description': 'Cloud Account Name',
          'type': 'string'
        },
        'gce_config': {
          'properties': {
            'ssh_pub_key': {
              'description': 'SSH Public Key',
              'type': 'string'
            },
            'zone': {
              'default': 'us-east1-b',
              'description': 'Zone',
              'type': 'string'
            }
          },
          'type': 'object'
        },
        'master_node_size': {
          'default': 'n1-standard-1',
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

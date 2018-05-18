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
          'description': 'The Supergiant cloud account you created for use with GCE.',
          'type': 'string'
        },
        'gce_config': {
          'properties': {
            'ssh_pub_key': {
              'description': 'The public key that will be used to SSH into the kube.',
              'type': 'string'
            },
            'zone': {
              'default': 'us-east1-b',
              'description': 'The GCE zone the kube will be created in.',
              'type': 'string'
            }
          },
          'type': 'object'
        },
        'master_node_size': {
          'default': 'n1-standard-1',
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

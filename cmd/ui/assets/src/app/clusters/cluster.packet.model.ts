export class ClusterPacketModel {
  packet = {
    'data': {
      'cloud_account_name': '',
      'master_node_size': 't1.small (Type 0)',
      'kube_master_count': 1,
      'ssh_pub_key': '',
      'name': '',
      'node_sizes': [
        't1.small (Type 0)',
        'c1.small (Type 1)',
        'x1.small (Type 1E)',
        'c2.medium.x86',
        'c1.large.arm (Type 2A)',
        's1.large (Type S)',
        'c1.xlarge (Type 3)',
        'm1.xlarge (Type 2)',
        'm2.xlarge.x86'
      ],
      'packet_config': {
        'facility': 'ewr1',
        'project': '',
      }
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
        'packet_config': {
          'properties': {
            'facility': {
              'default': 'ewr1',
              'description': 'The Packet facility (region) the kube will be created in.',
              'type': 'string'
            },
            'project': {
              'description': 'The Packet project the kube will be created in.',
              'type': 'string'
            }
          },
          'type': 'object'
        },
        'ssh_pub_key': {
          'description': 'The public key that will be used to SSH into the kube.',
          'type': 'string',
          'widget': 'textarea',
        },
        'master_node_size': {
          'default': 't1.small (Type 0)',
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
          'items': {
            'type': 'string'
          },
          'type': 'array'
        },
      }
    }
  };
}

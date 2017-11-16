export class ClusterPacketModel {
  packet = {
    'model': {
      'cloud_account_name': '',
      'master_node_size': 'Type 0',
      'kube_master_count': 1,
      'ssh_pub_key': '',
      'name': '',
      'node_sizes': [
        'Type 0',
        'Type 1',
        'Type 2',
        'Type 3',
        'Type 2A'
      ],
      'packet_config': {
        'facility': 'ewr1',
        'project': '',
      }
    },
    'schema': {
      'properties': {
        'cloud_account_name': {
          'description': 'Cloud Account Name',
          'type': 'string'
        },
        'master_node_size': {
          'default': 'Type 0',
          'description': 'Master Node Size',
          'type': 'string'
        },
        'kube_master_count': {
          'description': 'Kube Master Count',
          'type': 'number',
          'widget': 'number',
        },
        'ssh_pub_key': {
          'description': 'SSH Public Key',
          'type': 'string',
          'widget': 'textarea',
        },
        'name': {
          'description': 'Name',
          'type': 'string',
          'pattern': '^[a-z]([-a-z0-9]*[a-z0-9])?$',
          'maxLength': 12
        },
        'node_sizes': {
          'description': 'Node Sizes',
          'items': {
            'type': 'string'
          },
          'type': 'array'
        },
        'packet_config': {
          'properties': {
            'facility': {
              'default': 'ewr1',
              'description': 'Facility',
              'type': 'string'
            },
            'project': {
              'description': 'Project',
              'type': 'string'
            }
          },
          'type': 'object'
        }
      }
    }
  };
}

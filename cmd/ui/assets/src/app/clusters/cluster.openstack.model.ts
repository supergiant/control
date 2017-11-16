export class ClusterOpenStackModel {
  openstack = {
    'model': {
      'cloud_account_name': '',
      'master_node_size': 'm1.smaller',
      'name': '',
      'kube_master_count': 1,
      'node_sizes': [
        'm1.smaller',
        'm1.small'
      ],
      'openstack_config': {
        'image_name': 'CoreOS',
        'region': 'RegionOne',
        'public_gateway_id': '',
        'ssh_key_fingerprint': ''
      },
      'ssh_pub_key': ''
    },
    'schema': {
      'properties': {
        'cloud_account_name': {
          'description': 'Cloud Account Name',
          'type': 'string'
        },
        'master_node_size': {
          'default': 'm1.smaller',
          'description': 'Master Node Size',
          'type': 'string'
        },
        'name': {
          'description': 'Name',
          'type': 'string',
          'pattern': '^[a-z]([-a-z0-9]*[a-z0-9])?$',
          'maxLength': 12
        },
        'node_sizes': {
          'description': 'Node Sizes',
          'id': '/properties/node_sizes',
          'items': {
            'id': '/properties/node_sizes/items',
            'type': 'string'
          },
          'type': 'array'
        },
        'openstack_config': {
          'properties': {
            'image_name': {
              'default': 'CoreOS',
              'description': 'Image Name',
              'type': 'string'
            },
            'region': {
              'default': 'RegionOne',
              'description': 'Region',
              'type': 'string'
            },
            'public_gateway_id': {
              'description': 'public_gateway_id',
              'type': 'string'
            },
            'kube_master_count': {
              'description': 'Kube Master Count',
              'type': 'number',
              'widget': 'number',
            },
            'ssh_key_fingerprint': {
              'description': 'SSH Key Fingerprint',
              'type': 'string'
            }
          },
          'type': 'object'
        },
        'ssh_pub_key': {
          'description': 'SSH Public Key',
          'type': 'string'
        }
      }
    }
  };
}

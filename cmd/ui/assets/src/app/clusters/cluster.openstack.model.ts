export class ClusterOpenStackModel {
  openstack = {
    'data': {
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
          'description': 'The Supergiant cloud account you created for use with Openstack.',
          'type': 'string'
        },
        'openstack_config': {
          'properties': {
            'image_name': {
              'default': 'CoreOS',
              'description': 'The image the servers created will use.',
              'type': 'string'
            },
            'region': {
              'default': 'RegionOne',
              'description': 'The OpenStack region the kube will be created in.',
              'type': 'string'
            },
            'public_gateway_id': {
              'description': 'The gateway ID for your OpenStack public gateway.',
              'type': 'string'
            },
            'kube_master_count': {
              'description': 'The number of masters desired--for High Availability.',
              'type': 'number',
              'widget': 'number',
            },
            'ssh_key_fingerprint': {
              'description': 'The fingerprint of the public key that you uploaded to your OpenStack account.',
              'type': 'string'
            }
          },
          'type': 'object'
        },
        'ssh_pub_key': {
          'description': 'The public key that will be used to SSH into the kube.',
          'type': 'string'
        },
        'master_node_size': {
          'default': 'm1.smaller',
          'description': 'The size of the server the master will live on.',
          'type': 'string'
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

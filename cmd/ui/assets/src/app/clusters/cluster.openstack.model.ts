export class ClusterOpenStackModel {
  openstack = {
    'data': {
      'cloud_account_name': '',
      'master_node_size': 'm1.smaller',
      'name': '',
      'kube_master_count': 1,
      'kubernetes_version': '1.8.7',
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
          'type': 'string',
          'description': 'The desired name of the kube. Max length of 12 characters.',
          'pattern': '^[a-z]([-a-z0-9]*[a-z0-9])?$',
          'maxLength': 12
        },
        'kubernetes_version': {
          'type': 'string',
          'description': 'The Version of Kubernetes to be deployed.',
          'default': '1.5.7',
          'enum': ['1.5.7', '1.6.7', '1.7.7', '1.8.7', '1.9.7', '1.10.6', '1.11.1'] // TODO: <-- Should be dynamically populated.
        },
        'kube_master_count': { // TODO: this was originally in openstack_config, but only in schema?
          'type': 'number',
          'description': 'The number of masters desired - for High Availability.',
          'widget': 'number',
        },
        'cloud_account_name': {
          'type': 'string',
          'description': 'The Supergiant cloud account you created for use with Openstack.'
        },
        'openstack_config': {
          'type': 'object',
          'properties': {
            'image_name': {
              'type': 'string',
              'description': 'The image the servers created will use.',
              'default': 'CoreOS'
            },
            'region': {
              'type': 'string',
              'description': 'The OpenStack region the kube will be created in.',
              'default': 'RegionOne'
            },
            'public_gateway_id': {
              'type': 'string',
              'description': 'The gateway ID for your OpenStack public gateway.'
            },
            'ssh_key_fingerprint': {
              'type': 'string',
              'description': 'The fingerprint of the public key that you uploaded to your OpenStack account.'
            }
          }
        },
        'ssh_pub_key': {
          'type': 'string',
          'description': 'The public key that will be used to SSH into the kube.',
          'widget': 'textarea'
        },
        'master_node_size': {
          'type': 'string',
          'description': 'The size of the server the master will live on.',
          'default': 'm1.smaller'
        },
        'node_sizes': {
          'type': 'array',
          'description': 'The sizes you want to be available to Supergiant when scaling.',
          'id': '/properties/node_sizes',
          'items': {
            'id': '/properties/node_sizes/items',
            'type': 'string'
          }
        }
      }
    },
    'layout': [
      {
        'type': 'section',
        'title': 'Cluster Details',
        'items': [
          { 'key': 'name' },
          { 'key': 'kubernetes_version' },
          { 'key': 'cloud_account_name', 'readonly': true }
        ]
      },
      {
        'type': 'section',
        'title': 'Openstack Config',
        'items': [
          { 'key': 'openstack_config.image_name' },
          { 'key': 'openstack_config.region' },
          { 'key': 'openstack_config.public_gateway_id' },
          { 'key': 'openstack_config.ssh_key_fingerprint' }
        ]
      },
      {
        'type': 'section',
        'items': [
          { 'key': 'ssh_pub_key' },
          { 'key': 'kube_master_count' },
          { 'key': 'master_node_size' }
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

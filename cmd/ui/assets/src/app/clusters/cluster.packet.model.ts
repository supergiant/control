export class ClusterPacketModel {
  packet = {
    'data': {
      'cloud_account_name': '',
      'master_node_size': 't1.small (Type 0)',
      'kube_master_count': 1,
      'kubernetes_version': '1.8.7',
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
          'description': 'The Supergiant cloud account you created for use with Packet.'
        },
        'packet_config': {
          'type': 'object',
          'properties': {
            'facility': {
              'type': 'string',
              'description': 'The Packet facility (region) the kube will be created in.',
              'default': 'ewr1',
              'enum': [
                'sea1',
                'sjc1',
                'lax1',
                'dfw1',
                'ord1',
                'yyz1',
                'atl1',
                'iad1',
                'ewr1',
                'ams1',
                'fra1',
                'sin1',
                'hkg1',
                'nrt1',
                'syd1'
              ]
            },
            'project': {
              'type': 'string',
              'description': 'The Packet project the kube will be created in.'
            }
          },
          'required': [ 'facility' ]
        },
        'ssh_pub_key': {
          'type': 'string',
          'description': 'The public key that will be used to SSH into the kube.',
          'widget': 'textarea',
        },
        'master_node_size': {
          'type': 'string',
          'description': 'The size of the server the master will live on.',
          'default': 't1.small (Type 0)'
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
          'items': {
            'type': 'string'
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
        'title': 'Packet Config',
        'items': [
          { 'key': 'packet_config.facility' },
          { 'key': 'packet_config.project' }
        ]
      },
      {
        'type': 'section',
        'items': [
          { 'key': 'ssh_pub_key' },
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

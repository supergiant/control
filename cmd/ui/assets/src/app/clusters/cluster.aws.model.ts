export class ClusterAWSModel {
  aws = {
    'data': {
      'name': '',
      'aws_config': {
        'region': 'us-east-1',
        'vpc_ip_range': '172.20.0.0/16'
      },
      'cloud_account_name': '',
      'master_node_size': 'm4.large',
      'ssh_pub_key': '',
      'kube_master_count': 1,
      'kubernetes_version': '1.8.7',
      'node_sizes': [
        'm4.large',
        'm4.xlarge',
        'm4.2xlarge',
        'm4.4xlarge'
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
          'description': 'The Supergiant cloud account you created for use with AWS.',
          'type': 'string'
        },
        'aws_config': {
          'properties': {
            'region': {
              'default': 'us-east-1',
              'description': 'The AWS region the kube will be created in.',
              'type': 'string'
            },
            'vpc_ip_range': {
              'default': '172.20.0.0/16',
              'description': 'The range of IP addresses you want available to the kube.',
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
          'default': 'm4.large',
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
          'widget': 'array',
          'items': {
            'type': 'string'
          },
          'type': 'array'
        }
      }
    },
    'layout': [
      { 'type': 'flex', 'flex-flow': 'row wrap', 'items': ['name', 'region'] },
    ]
  };
}

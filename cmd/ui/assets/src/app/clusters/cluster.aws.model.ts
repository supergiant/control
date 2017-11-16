export class ClusterAWSModel {
  aws = {
    'model': {
      'aws_config': {
        'region': 'us-east-1',
        'vpc_ip_range': '172.20.0.0/16'
      },
      'cloud_account_name': '',
      'master_node_size': 'm4.large',
      'name': '',
      'ssh_pub_key': '',
      'kube_master_count': 1,
      'node_sizes': [
        'm4.large',
        'm4.xlarge',
        'm4.2xlarge',
        'm4.4xlarge'
      ]
    },
    'schema': {
      'properties': {
        'aws_config': {
          'properties': {
            'region': {
              'default': 'us-east-1',
              'description': 'Region',
              'type': 'string'
            },
            'vpc_ip_range': {
              'default': '172.20.0.0/16',
              'description': 'VPC IP Range',
              'type': 'string'
            }
          },
          'type': 'object'
        },
        'cloud_account_name': {
          'description': 'Cloud Account Name',
          'type': 'string'
        },
        'master_node_size': {
          'default': 'm4.large',
          'description': 'Master Node Size',
          'type': 'string'
        },
        'name': {
          'description': 'Name (a-z,0-9)',
          'type': 'string',
          'pattern': '^[a-z]([-a-z0-9]*[a-z0-9])?$',
          'maxLength': 12
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
        'node_sizes': {
          'description': 'Node Sizes',
          'widget': 'array',
          'items': {
            'type': 'string'
          },
          'type': 'array'
        }
      }
    }
  };
}

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
          'type': 'string',
          'description': 'The desired name of the kube. Max length of 12 characters.',
          'pattern': '^[a-z]([-a-z0-9]*[a-z0-9])?$',
          'maxLength': 12
        },
        'kubernetes_version': {
          'type': 'string',
          'default': '1.5.7',
          'description': 'The Version of Kubernetes to be deployed.',
          'enum': ['1.5.7', '1.6.7', '1.7.7', '1.8.7'] // TODO: <-- Should be dynamically populated.
        },
        'cloud_account_name': { // TODO: <-- Should be dynamically populated enum with any cloud accounts linked to provider
          'type': 'string',
          'description': 'The Supergiant cloud account you created for use with AWS.'
        },
        'aws_config': {
          'type': 'object',
          'properties': {
            'region': {
              'type': 'string',
              'description': 'The AWS region the kube will be created in.',
              'default': 'us-east-1',
              'enum': [
                'us-east-1',
                'us-east-2',
                'us-west-1',
                'us-west-2',
                'ap-northeast-1',
                'ap-northeast-2',
                'ap-northeast-3',
                'ap-south-1',
                'ap-southeast-1',
                'ap-southeast-2',
                'ca-central-1',
                'cn-north-1',
                'cn-northwest-1',
                'eu-central-1',
                'eu-west-1',
                'eu-west-2',
                'eu-west-3',
                'sa-east-1'
              ] // TODO: <-- Should be dynamically populated.
            },
            'vpc_ip_range': {
              'default': '172.20.0.0/16',
              'description': 'The range of IP addresses you want available to the kube.',
              'type': 'string'
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
          'default': 'm4.large',
          'description': 'The size of the server the master will live on.'
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
          'widget': 'array',
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
        'title': 'AWS Config',
        'items': [
          { 'key': 'aws_config.region' },
          { 'key': 'aws_config.vpc_ip_range' }
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

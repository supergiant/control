export class PodsModel {
  pod = {
    'model': {
      'kind': 'Pod',
      'kube_name': '',
      'name': '',
      'namespace': 'default',
      'resource': {
        'metadata': {
          'labels': {}
        },
        'spec': {
          'containers': [
            {
              'image': 'jenkins',
              'name': 'jenkins',
              'ports': [
                {
                  'containerPort': 8080
                }
              ]
            }
          ]
        }
      }
    },
    'schema': {
      'properties': {
        'kind': {
          'default': 'Pod',
          'description': 'Kind',
          'type': 'string'
        },
        'kube_name': {
          'description': 'Kube Name',
          'type': 'string'
        },
        'name': {
          'description': 'Name',
          'type': 'string'
        },
        'namespace': {
          'default': 'default',
          'description': 'Namespace',
          'type': 'string'
        },
        'resource': {
          'properties': {
            'metadata': {
              'properties': {
                'labels': {
                  'description': 'Labels',
                  'type': 'string'
                }
              },
              'type': 'object'
            },
            'spec': {
              'properties': {
                'containers': {
                  'description': 'Containers',
                  'type': 'string'
                }
              },
              'type': 'object'
            }
          },
          'type': 'object'
        }
      }
    }
  };
  public providers = {
    'pod': this.pod
  };
}

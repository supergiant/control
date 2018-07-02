export class KubeResourcesModel {
  pod = {
    'model': {
      'kind': 'Pod',
      'kube_name': '',
      'name': '',
      'namespace': 'default',
      'resource': {
        'metadata': {
          'labels': []
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
      'definitions': {
        'containers_obj': {
          'type': 'object',
          'properties': {
            'image': { 'type': 'string' },
            'name': { 'type': 'string' },
            'ports': {
              'type': 'array',
              'items': { '$ref': '#/definitions/ports_obj' }
            }
          }
        },
        'ports_obj': {
          'type': 'object',
          'properties': {
            'containerPort': { 'type': 'number'}
          }
        }
      },
      'properties': {
        'name': {
          'type': 'string',
          'description': 'Name of resource'
        },
        'kube_name': {
          'type': 'string',
          'description': 'The Kube this resource will be attached to.'
        },
        'namespace': {
          'type': 'string',
          'description': 'Namespace the resource will be under.',
          'default': 'default'
        },
        'kind': {
          'type': 'string',
          'description': 'Type of resource to create',
          'default': 'Pod',
          'readonly': true,
          'htmlClass': 'hideMe'
        },
        'resource': {
          'type': 'object',
          'notitle': true,
          'properties': {
            'metadata': {
              'type': 'object',
              'notitle': true,
              'properties': {
                'labels': {
                  'type': 'array',
                  'items': {
                    'type': 'object',
                    'properties': {
                      'key': { 'type': 'string' },
                      'value': { 'type': 'string' }
                    }
                  }
                }
              }
            },
            'spec': {
              'type': 'object',
              'notitle': true,
              'properties': {
                'containers': {
                  'type': 'array',
                  'title': 'Containers',
                  'items': { '$ref': '#/definitions/containers_obj', 'extendRefs': true }
                }
              }
            }
          }
        }
      }
    },
    'layout' : [
      // angular2-json-schema-form doesn't support custom layouts with '$ref'-ed nested arrays
      // https://github.com/dschnelldavis/angular2-json-schema-form/issues/149
      { 'type': 'section', 'title': 'Pod Details' },
      '*',
      { 'type': 'submit', 'title': 'Create' }
    ]
  };

  service = {
    'model': {
      'kind': 'Service',
      'kube_name': '',
      'name': '',
      'namespace': 'default',
      'template': {
        'spec': {
          'ports': [
            { 'name': 'jenkins', 'port': 8080 }
          ],
          'selector': [],
          'type': 'NodePort'
        }
      }
    },
    'schema': {
      'properties': {
        'kind': {
          'type': 'string',
          'description': 'Type of resource to create',
          'default': 'Service',
        },
        'kube_name': {
          'type': 'string',
          'description': 'The Kube this resource will be attached to.'
        },
        'name': {
          'type': 'string',
          'description': 'Name of resource'
        },
        'namespace': {
          'type': 'string',
          'description': 'Namespace the resource will be under.',
          'default': 'default'
        },
        'template': {
          'type': 'object',
          'properties': {
            'spec': {
              'type': 'object',
              'properties': {
                'ports': {
                  'type': 'array',
                    'items': {
                      'type': 'object',
                      'properties': {
                        'name': { 'type': 'string' },
                        'port': { 'type': 'number' }
                      }
                  }
                },
                'selector': {
                  'type': 'array',
                    'items': {
                      'type': 'object',
                      'properties': {
                        'key': { 'type': 'string' },
                        'value': { 'type': 'string' }
                    }
                  }
                },
                'type': {
                  'type': 'string',
                  'description': 'Type',
                  'default': 'NodePort',
                }
              }
            }
          }
        }
      }
    },
    'layout': [
      { 'type': 'section',
        'title': 'Service Details',
        'items': [
          { 'key': 'name' },
          { 'key': 'kube_name', 'readonly': true },
          { 'key': 'namespace' },
        ]
      },
      { 'type': 'array',
        'title': 'Ports',
        'display': 'flex',
        'flex-direction': 'column',
        'items': [
          { 'type': 'flex',
            'flex-flow': 'row-wrap',
            'items': [
              { 'key': 'template.spec.ports[].name' },
              { 'key': 'template.spec.ports[].port' }
            ]
          }
        ]
      },
      { 'type': 'array',
        'title': 'Selectors',
        'display': 'flex',
        'flex-direction': 'column',
        'items': [
          { 'type': 'flex',
            'flex-flow': 'row-wrap',
            'items': [
              { 'key': 'template.spec.selector[].key' },
              { 'key': 'template.spec.selector[].value' }
            ]
          }
        ]
      },
      { 'key': 'template.spec.type' },
      { 'type': 'submit', 'title': 'Create' }
    ]
  };

  loadBalancer = {
    'model': {
      'kind': 'LoadBalancer',
      'kube_name': '',
      'name': '',
      'namespace': 'default',
      'ports': [
        { 'key': '80', 'value': 8080 }
      ],
      'selector': []
    },
    'schema': {
      'properties': {
        'kind': {
          'type': 'string',
          'description': 'Type of resource to create',
          'default': 'LoadBalancer'
        },
        'kube_name': {
          'type': 'string',
          'description': 'The Kube this resource will be attached to.'
        },
        'name': {
          'type': 'string',
          'description': 'Name'
        },
        'namespace': {
          'type': 'string',
          'description': 'Namespace the resource will be under.',
          'default': 'default'
        },
        'ports': {
          'type': 'array',
            'items': {
              'type': 'object',
              'properties': {
                'key': { 'type': 'string' },
                'value': { 'type': 'number' }
              }
          }
        },
        'selector': {
          'type': 'array',
            'items': {
              'type': 'object',
              'properties': {
                'key': { 'type': 'string' },
                'value': { 'type': 'string' }
            }
          }
        }
      }
    },
    'layout': [
      { 'type': 'section',
        'title': 'LoadBalancer Details',
        'items': [
          { 'key': 'name' },
          { 'key': 'kube_name', 'readonly': true },
          { 'key': 'namespace' },
        ]
      },
      { 'type': 'array',
        'title': 'Ports',
        'display': 'flex',
        'flex-direction': 'column',
        'items': [
          { 'type': 'flex',
            'flex-flow': 'row-wrap',
            'items': [
              { 'key': 'ports[].key', 'title': 'Name' },
              { 'key': 'ports[].value', 'title': 'Port' }
            ]
          }
        ]
      },
      { 'type': 'array',
        'title': 'Selectors',
        'display': 'flex',
        'flex-direction': 'column',
        'items': [
          { 'type': 'flex',
            'flex-flow': 'row-wrap',
            'items': [
              { 'key': 'selector[].key' },
              { 'key': 'selector[].value' }
            ]
          }
        ]
      },
      { 'type': 'submit', 'title': 'Create' }
    ]
  };

  public providers = {
    pod: this.pod,
    service: this.service,
    loadBalancer: this.loadBalancer
  }
}

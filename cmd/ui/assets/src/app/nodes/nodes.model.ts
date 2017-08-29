export class NodesModel {
  node = {
    'model': {
      'kube_name': '',
      'size': '',
    },
    'schema': {
      'properties': {
        'kube_name': {
          'description': 'Kube Name',
          'type': 'string'
        },
        'size': {
          'description': 'Size',
          'type': 'string',
          'widget': 'select',
          'oneOf': [],
          'default': '',
        }
      }
    }
  };
  public providers = {
    'node': this.node
  };
}

export class AppsModel {
  app = {
    'model': {
      'chart_name': '',
      'chart_version': '',
      'config': null,
      'kube_name': '',
      'name': '',
      'repo_name': '',
      'namespace': ''
    },
    'schema': {}
  };
  public providers = {
    'app': this.app
  };
}

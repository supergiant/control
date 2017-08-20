export class AppsModel {
  app = {
    'model': {
      'chart_name': '',
      'chart_version': '',
      'config': null,
      'kube_name': '',
      'name': '',
      'repo_name': ''
    },
    'schema': {}
  };
  public providers = {
    'app': this.app
  };
}

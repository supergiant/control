export class AppsModel {
  app = {
    'model': {
      'chart_name': '',
      'chart_version': '',
      'config': null,
      // this needs to be called cluster in display
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

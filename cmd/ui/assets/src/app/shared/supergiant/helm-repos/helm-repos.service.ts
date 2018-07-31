import { Injectable } from '@angular/core';
import { UtilService } from '../util/util.service';

@Injectable()
export class HelmRepos {
  helmReposPath = '/v1/api/helm_repos';

  constructor(private util: UtilService) { }
  public get(id?) {
    if (id) {
      return this.util.fetch(this.helmReposPath + '/' + id);
    }
    return this.util.fetch(this.helmReposPath);
  }
  public create(data) {
    return this.util.post(this.helmReposPath, data);
  }
  public update(id, data) {
    return this.util.update(this.helmReposPath + '/' + id, data);
  }
  public delete(id) {
    return this.util.destroy(this.helmReposPath + '/' + id);
  }
}

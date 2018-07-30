import { Injectable } from '@angular/core';
import { UtilService } from '../util/util.service';

@Injectable()
export class KubeResources {
  kubeResourcesPath = '/v1/api/kube_resources';

  constructor(private util: UtilService) { }
  public get(id?) {
    if (id) {
      return this.util.fetch(this.kubeResourcesPath + '/' + id);
    }
    return this.util.fetch(this.kubeResourcesPath);
  }
  public create(data) {
    return this.util.post(this.kubeResourcesPath, data);
  }
  public start(id, data) {
    return this.util.post(this.kubeResourcesPath + '/' + id + '/start', data);
  }
  public stop(id, data) {
    return this.util.post(this.kubeResourcesPath + '/' + id + '/stop', data);
  }
  public update(id, data) {
    return this.util.update(this.kubeResourcesPath + '/' + id, data);
  }
  public delete(id) {
    return this.util.destroy(this.kubeResourcesPath + '/' + id);
  }
}

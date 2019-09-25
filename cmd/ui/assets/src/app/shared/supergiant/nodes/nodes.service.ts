import { Injectable } from '@angular/core';
import { UtilService } from '../util/util.service';

@Injectable()
export class Nodes {
  nodesPath = '/api/v1/nodes';

  constructor(private util: UtilService) { }
  public get(id?) {
    if (id) {
      return this.util.fetch(this.nodesPath + '/' + id);
    }
    return this.util.fetch(this.nodesPath);
  }
  public create(data) {
    return this.util.post(this.nodesPath, data);
  }
  public update(id, data) {
    return this.util.update(this.nodesPath + '/' + id, data);
  }
  public delete(kubeId, nodeId) {
    // FIXME remove kube id from backand requests, it's redundant
    return this.util.destroy(`/api/v1/kubes/${kubeId}/nodes/${nodeId}`);
  }
}

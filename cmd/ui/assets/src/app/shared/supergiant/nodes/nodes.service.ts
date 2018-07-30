import { Injectable } from '@angular/core';
import { UtilService } from '../util/util.service';

@Injectable()
export class Nodes {
  nodesPath = '/v1/api/nodes';

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
  public delete(id) {
    return this.util.destroy(this.nodesPath + '/' + id);
  }
}

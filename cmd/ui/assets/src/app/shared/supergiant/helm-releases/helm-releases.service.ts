import { Injectable } from '@angular/core';
import {UtilService} from '../util/util.service'

@Injectable()
export class HelmReleases{
  helmReleasesPath = "/api/v0/helm_releases"

  constructor(private util: UtilService) {}
  public get(id?){
    if (id) {
      this.util.fetch(this.helmReleasesPath +"/" + id)
    }
    return this.util.fetch(this.helmReleasesPath)
  }
  public create(data) {
    return this.util.post(this.helmReleasesPath, data)
  }
  public update(id, data) {
    return this.util.update(this.helmReleasesPath+ "/" + id, data)
  }
  public delete(id){
    return this.util.destroy(this.helmReleasesPath + "/" + id)
  }
}

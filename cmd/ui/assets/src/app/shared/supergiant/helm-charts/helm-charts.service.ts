import { Injectable } from '@angular/core';
import {UtilService} from '../util/util.service'

@Injectable()
export class HelmCharts{
  helmChartsPath = "/api/v0/helm_charts"

  constructor(private util: UtilService) {}
  public get(id?){
    if (id) {
      this.util.fetch(this.helmChartsPath +"/" + id)
    }
    return this.util.fetch(this.helmChartsPath)
  }
  public create(data) {
    return this.util.post(this.helmChartsPath, data)
  }
  public update(id, data) {
    return this.util.update(this.helmChartsPath+ "/" + id, data)
  }
  public delete(id){
    return this.util.destroy(this.helmChartsPath + "/" + id)
  }
}

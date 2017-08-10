import { Injectable } from '@angular/core';
import {UtilService} from '../util/util.service'

@Injectable()
export class LoadBalancers{
  loadBalancersPath = "/api/v0/load_balancers"

  constructor(private util: UtilService) {}
  public get(id?){
    if (id) {
      this.util.fetch(this.loadBalancersPath +"/" + id)
    }
    return this.util.fetch(this.loadBalancersPath)
  }
  public schema(){
    return this.util.fetch(this.loadBalancersPath + "/schema")
  }
  public create(data) {
    return this.util.post(this.loadBalancersPath, data)
  }
  public update(id, data) {
    return this.util.update(this.loadBalancersPath+ "/" + id, data)
  }
  public delete(id){
    return this.util.destroy(this.loadBalancersPath + "/" + id)
  }
}

import { Injectable } from '@angular/core';
import {UtilService} from '../util/util.service'

@Injectable()
export class Kubes{
  kubesPath = "/api/v0/kubes"

  constructor(private util: UtilService) {}
  public get(id?){
    if (id) {
      this.util.fetch(this.kubesPath +"/" + id)
    }
    return this.util.fetch(this.kubesPath)
  }
  public schema(){
    return this.util.fetch(this.kubesPath + "/schema")
  }
  public create(data) {
    return this.util.post(this.kubesPath, data)
  }
  public provision(id, data) {
    return this.util.post(this.kubesPath+ "/" + id + "/provision", data)
  }
  public update(id, data) {
    return this.util.update(this.kubesPath+ "/" + id, data)
  }
  public delete(id){
    return this.util.destroy(this.kubesPath + "/" + id)
  }
}

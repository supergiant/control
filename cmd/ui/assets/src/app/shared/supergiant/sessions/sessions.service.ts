import { Injectable } from '@angular/core';
import {UtilService} from '../util/util.service'

@Injectable()
export class Sessions{
  sessionsPath = "/api/v0/sessions"

  constructor(private util: UtilService) {}
  public get(id?){
    if (id) {
      this.util.fetch(this.sessionsPath +"/" + id)
    }
    return this.util.fetch(this.sessionsPath)
  }
  public create(data) {
    return this.util.post(this.sessionsPath, data)
  }
  public delete(id){
    return this.util.destroy(this.sessionsPath + "/" + id)
  }
}

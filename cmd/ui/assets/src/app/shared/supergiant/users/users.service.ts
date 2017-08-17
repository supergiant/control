import { Injectable } from '@angular/core';
import {UtilService} from '../util/util.service'

@Injectable()
export class Users{
  usersPath = "/api/v0/users"

  constructor(private util: UtilService) {}
  public get(id?){
    if (id) {
      this.util.fetch(this.usersPath +"/" + id)
    }
    return this.util.fetch(this.usersPath)
  }
  public create(data) {
    return this.util.post(this.usersPath, data)
  }
  public update(id, data) {
    return this.util.update(this.usersPath+ "/" + id, data)
  }
  public delete(id){
    return this.util.destroy(this.usersPath + "/" + id)
  }
  public generateToken(id){
    return this.util.post(this.usersPath + "/" + id + "/regenerate_api_token", "")
  }
}

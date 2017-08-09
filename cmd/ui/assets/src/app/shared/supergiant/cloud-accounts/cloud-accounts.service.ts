import { Injectable } from '@angular/core';
import {UtilService} from '../util/util.service'

Injectable()
export class CloudAccount{
  cloudAccountsPath = "/api/v0/cloud_accounts"

  constructor(private util: UtilService) {}
  // Get Cloud Accounts object.
  public get(id?){
    if (id) {
      this.util.fetch(this.cloudAccountsPath +"/" + id)
    }
    return this.util.fetch(this.cloudAccountsPath)

  }
  public schema(){
    return this.util.fetch(this.cloudAccountsPath + "/schema")
  }
  // Creates a  Cloud Account object, requires a Cloud Account object.
  public create(data) {
    return this.util.post(this.cloudAccountsPath, data)
  }
  // Edit a Cloud Account object, requires a new CloudAccount object.
  public update(id, data) {
    return this.util.update(this.cloudAccountsPath+ "/" + id, data)
  }
  // Deletes a Cloud Account
  public delete(id){
    return this.util.destroy(this.cloudAccountsPath + "/" + id)
  }
}

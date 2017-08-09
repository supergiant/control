import { Injectable } from '@angular/core';
import { Http, Response, Headers } from '@angular/http';

export namespace Supergiant {
  export namespace CloudAccounts {
    var cloudAccountsPath = "/api/v0/cloud_accounts"
    // constructor(private util: Util) {}
    // Get Cloud Accounts object.
    export function get(id?){
      if (id) {
      return this.util.get(cloudAccountsPath +"/" + id)
      }
      return this.util.get(cloudAccountsPath)

    }
    export function schema(){
      return this.util.get(cloudAccountsPath + "/schema")
    }
    // Creates a  Cloud Account object, requires a Cloud Account object.
    export function create(data) {
      return this.util.post(cloudAccountsPath, data)
    }
    // Edit a Cloud Account object, requires a new CloudAccount object.
    export function update(id, data) {
      return this.util.update(cloudAccountsPath+ "/" + id, data)
    }
    // Deletes a Cloud Account
    export function destroy(id){
      return this.util.delete(cloudAccountsPath + "/" + id)
    }
  }

export class Util {
  serverEndpoint = "http://localhost:8080"
  apiToken = 'SGAPI token="iPIrhzlIHfRtHkRXW8NoAou0pUGFUfmo"'

   constructor(private http: Http) {}
   get(path) {
    let headers = new Headers();
    headers.append('Authorization', this.apiToken);
    return this.http.get(this.serverEndpoint + path, { headers: headers })
  }

  post(path, data) {
    var json = JSON.stringify(data)
    let headers = new Headers();
    headers.append('Authorization', this.apiToken);
    return this.http.post(this.serverEndpoint + path, json, { headers: headers })
  }

  update(path, data) {
    var json = JSON.stringify(data)
    let headers = new Headers();
    headers.append('Authorization', this.apiToken);
    return this.http.put(this.serverEndpoint + path, json, { headers: headers })
  }

  delete(path) {
    let headers = new Headers();
    headers.append('Authorization', this.apiToken);
    return this.http.delete(this.serverEndpoint + path, { headers: headers })
  }
}
}

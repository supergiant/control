import { Injectable } from '@angular/core';
import { Http, Response, Headers } from '@angular/http';
import {Observable} from "RxJS/Rx";

@Injectable()
export class UtilService {
  serverEndpoint = "http://localhost:8080"
  apiToken = 'SGAPI token="3roIJ9WKhS0NLk1LBKalNgETfFxyrQSX"'

  constructor(private http: Http){}

  fetch(path) {
    let headers = new Headers();
    headers.append('Authorization', this.apiToken);
    return this.http.get(this.serverEndpoint + path, { headers: headers }).map(response => response.json())
  }

  post(path, data) {
    var json = JSON.stringify(data)
    let headers = new Headers();
    headers.append('Authorization', this.apiToken);
    return this.http.post(this.serverEndpoint + path, json, { headers: headers }).map(response => response.json())
  }

  update(path, data) {
    var json = JSON.stringify(data)
    let headers = new Headers();
    headers.append('Authorization', this.apiToken);
    return this.http.put(this.serverEndpoint + path, json, { headers: headers }).map(response => response.json())
  }

  destroy(path) {
    let headers = new Headers();
    headers.append('Authorization', this.apiToken);
    return this.http.delete(this.serverEndpoint + path, { headers: headers }).map(response => response.json())
  }
}

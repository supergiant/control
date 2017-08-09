import { Injectable } from '@angular/core';
import { Http, Response, Headers } from '@angular/http';

@Injectable()
export class UtilService {
  serverEndpoint = "http://localhost:8080"
  apiToken = 'SGAPI token="DM5ijdza8wWhkMcN7RRDmQeuziogI2C7"'

  constructor(private http: Http){}

  fetch(path) {
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

  destroy(path) {
    let headers = new Headers();
    headers.append('Authorization', this.apiToken);
    return this.http.delete(this.serverEndpoint + path, { headers: headers })
  }
}

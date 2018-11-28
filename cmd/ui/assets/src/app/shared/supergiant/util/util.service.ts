import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

/*
* @deprecated
* For dev server use proxy.conf.json
* For production use own config of the remote server
*/
@Injectable()
export class UtilService {
  serverEndpoint = '';

  constructor(
    private http: HttpClient,
  ) {
  }

  fetch(path) {
    return this.http.get<any>(this.serverEndpoint + path + '?limit=1000');
  }

  fetchResponse(path) {
    return this.http.get(this.serverEndpoint + path, { observe: "response" });
  }

  post(path, data) {
    const json = JSON.stringify(data);
    return this.http.post(this.serverEndpoint + path, json);
  }

  postResponse(path, data) {
    const json = JSON.stringify(data);
    return this.http.post(this.serverEndpoint + path, json, { observe: "response" });
  }

  update(path, data) {
    const json = JSON.stringify(data);
    return this.http.put(this.serverEndpoint + path, json);
  }

  destroy(path) {
    return this.http.delete(this.serverEndpoint + path, { headers: {} });
  }
}

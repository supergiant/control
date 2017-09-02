import { Injectable } from '@angular/core';
import { Http, Response, Headers } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import { Location } from '@angular/common';

@Injectable()
export class UtilService {
  serverEndpoint = 'http://localhost:8080';
  sessionToken: string;
  SessionID: string;

  constructor(
    private http: Http,
    private location: Location,
  ) {
    console.log('foo' + location.prepareExternalUrl(location.path()));
    console.log('window.location', window.location);
    console.log('window.location.href', window.location.href);
    console.log('window.location.origin', window.location.origin);
    this.serverEndpoint = window.location.origin + '/api';
  }

  fetch(path) {
    const headers = new Headers();
    headers.append('Authorization', this.sessionToken);
    return this.http.get(this.serverEndpoint + path + '?limit=1000', { headers: headers }).map(response => response.json());
  }

  fetchNoMap(path) {
    const headers = new Headers();
    headers.append('Authorization', this.sessionToken);
    return this.http.get(this.serverEndpoint + path, { headers: headers }).map(response => response);
  }

  post(path, data) {
    const json = JSON.stringify(data);
    const headers = new Headers();
    headers.append('Authorization', this.sessionToken);
    return this.http.post(this.serverEndpoint + path, json, { headers: headers }).map(response => response.json());
  }

  update(path, data) {
    const json = JSON.stringify(data);
    const headers = new Headers();
    headers.append('Authorization', this.sessionToken);
    return this.http.put(this.serverEndpoint + path, json, { headers: headers }).map(response => response.json());
  }

  destroy(path) {
    const headers = new Headers();
    headers.append('Authorization', this.sessionToken);
    return this.http.delete(this.serverEndpoint + path, { headers: headers }).map(response => response.json());
  }
}

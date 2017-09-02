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
    console.log('window.location.origin', window.location.origin);
    console.log('window.location.host', window.location.hostname);
    console.log('window.location.splice', window.location.pathname.split('/').splice(1).join('/'));
    console.log('window.location.split', window.location.pathname.split('/'));
    console.log('window.location.path', window.location.pathname);
    console.log('window.location.href', window.location.href);

    if (window.location.pathname.split('/')[2] === 'ui') {
      this.serverEndpoint = '/' + window.location.pathname.split('/')[1] + '/server';
    } else {
      this.serverEndpoint = 'http://localhost:8080';
    }
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

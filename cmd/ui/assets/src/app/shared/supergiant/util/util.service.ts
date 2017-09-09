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
    if (window.location.pathname.split('/')[2] === 'ui') {
      this.serverEndpoint = '/' + window.location.pathname.split('/')[1] + '/server';
    } else {
      console.log(window.location.protocol);
      if (window.location.hostname === 'localhost') {
        this.serverEndpoint = window.location.protocol + '//' + window.location.hostname + '8080';
      } else {
        this.serverEndpoint = window.location.protocol + '//' + window.location.hostname;
      }
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

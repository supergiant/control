
import {throwError as observableThrowError,  Observable } from 'rxjs';

import {catchError} from 'rxjs/operators';
import { Injectable } from '@angular/core';
// REFACTOR: remove deprecated modules
import { Request, XHRBackend, RequestOptions, Response, Http, RequestOptionsArgs, Headers } from '@angular/http';


import { Router } from '@angular/router';


@Injectable()
export class AuthenticatedHttpService extends Http {

  constructor(backend: XHRBackend, defaultOptions: RequestOptions, private router: Router, ) {
    super(backend, defaultOptions);
  }

  // REFACTOR: remove deprecated method
  request(url: string | Request, options?: RequestOptionsArgs): Observable<Response> {
    return super.request(url, options).pipe(catchError((error: Response) => {
      if ((error.status === 401 || error.status === 403) && (window.location.href.match(/\?/g) || []).length < 2) {
        console.log('User session has expired... Redirect to login.');
        this.router.navigate(['']);
        return observableThrowError(null);
      }
      return observableThrowError(error);
    }));
  }
}

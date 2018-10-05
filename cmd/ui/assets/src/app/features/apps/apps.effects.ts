import { Injectable } from '@angular/core';
import { Actions, Effect, ofType } from '@ngrx/effects';
import { AppsActionTypes, SetCharts } from './apps.actions';
import { map, switchMap } from 'rxjs/operators';
import { HttpClient } from '@angular/common/http';
import { Action } from '@ngrx/store';
import { Observable } from 'rxjs';

@Injectable()
export class AppsEffects {

  @Effect()
  loadFoos$: Observable<Action> = this.actions$.pipe(
    ofType(AppsActionTypes.LoadAppss),
    switchMap(() => this.http.get('/v1/api/helm/repositories/supergiant/charts')),
    map(charts => new SetCharts(charts)),
  );

  constructor(
    private actions$: Actions,
    private http: HttpClient,
  ) {
  }
}

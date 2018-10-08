import { Injectable }                                       from '@angular/core';
import { Actions, Effect, ofType }                          from '@ngrx/effects';
import { AppStoreActionTypes, LoadSupergiantChartsSuccess } from './apps.actions';
import { map, switchMap }                                   from 'rxjs/operators';
import { HttpClient }                                       from '@angular/common/http';
import { Action }                                           from '@ngrx/store';
import { Observable }                                       from 'rxjs';
import * as appsActions                                     from './actions'

@Injectable()
export class AppsEffects {

  @Effect()
  loadSupergiantCharts: Observable<Action> = this.actions$.pipe(
    ofType(AppStoreActionTypes.LoadSupergiantCharts),
    switchMap(() => this.http.get('/v1/api/helm/repositories/supergiant/charts')),
    map(charts => new LoadSupergiantChartsSuccess(charts)),
  );

  loadVerifiedCharts: Observable<Action> = this.actions$.pipe(
    ofType(appsActions.VerifiedAppActionTypes.LoadVerifiedApps),
    switchMap(() => this.http.get('/v1/api/helm/repositories/supergiant/charts')),
    map(charts => new appsActions.LoadVerifiedAppsSuccess(charts)),
  );

  constructor(
    private actions$: Actions,
    private http: HttpClient,
  ) {
  }
}

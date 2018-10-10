import { Injectable }                                       from '@angular/core';
import { Actions, Effect, ofType }                          from '@ngrx/effects';
import { AppStoreActionTypes, LoadSupergiantChartsSuccess } from './actions/supergiant-app-actions';
import { map, switchMap }                                   from 'rxjs/operators';
import { HttpClient }                                       from '@angular/common/http';
import { Action }                                           from '@ngrx/store';
import { Observable }                                       from 'rxjs';
import { LoadOtherAppsSuccess, OtherAppActionTypes }        from "./actions/other-app.actions";
import { Chart }                                            from "./apps.reducer";
import {
  AppDetailActions,
  AppDetailActionTypes, LoadAppDetails,
  LoadAppDetailsSuccess,
  LoadVerifiedAppsSuccess,
  VerifiedAppActionTypes
}                                                           from "./actions";


@Injectable()
export class AppsEffects {

  @Effect()
  loadSupergiantCharts: Observable<Action> = this.actions$.pipe(
    ofType(AppStoreActionTypes.LoadSupergiantCharts),
    switchMap(() => this.http.get('/v1/api/helm/repositories/supergiant/charts')),
    map(charts => new LoadSupergiantChartsSuccess(charts)),
  );

  @Effect()
  loadVerifiedCharts: Observable<Action> = this.actions$.pipe(
    ofType(VerifiedAppActionTypes.LoadVerifiedApps),
    switchMap(() => this.http.get('/v1/api/helm/repositories/supergiant/charts')),
    map((charts: Chart[]) => new LoadVerifiedAppsSuccess(charts)),
  );

  @Effect()
  loadOtherCharts: Observable<Action> = this.actions$.pipe(
    ofType(OtherAppActionTypes.LoadOtherApps),
    switchMap(() => this.http.get('/v1/api/helm/repositories/supergiant/charts')),
    map((charts: Chart[]) => new LoadOtherAppsSuccess(charts)),
  );

  @Effect()
  loadChartDetails: Observable<LoadAppDetailsSuccess> = this.actions$.pipe(
    ofType(AppDetailActionTypes.LoadAppDetails),
    map((action: LoadAppDetails) => action.payload),
    switchMap(({ repo, chart }) => this.http.get(`/v1/api/helm/repositories/${repo}/charts/${chart}`)),
    map((chart: Chart) => new LoadAppDetailsSuccess(chart)),
  );


  constructor(
    private actions$: Actions,
    private http: HttpClient,
  ) {
  }
}

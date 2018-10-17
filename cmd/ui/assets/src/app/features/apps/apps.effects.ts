import { Injectable }                                             from '@angular/core';
import { Actions, Effect, ofType }                                from '@ngrx/effects';
import { AppStoreActionTypes, LoadSupergiantChartsSuccess }       from './actions/supergiant-app-actions';
import { distinctUntilChanged, filter, map, mergeMap, switchMap } from 'rxjs/operators';
import { HttpClient }                                             from '@angular/common/http';
import { Action }                                                 from '@ngrx/store';
import { Observable }                                             from 'rxjs';
import { LoadOtherAppsSuccess, OtherAppActionTypes }              from "./actions/other-app.actions";
import { Chart }                                                  from "./apps.reducer";
import {
  AppCommonActions,
  AppCommonActionTypes,
  AppDetailActions,
  AppDetailActionTypes, LoadAppDetails,
  LoadAppDetailsSuccess, LoadCharts, LoadChartsSuccess,
  LoadVerifiedAppsSuccess,
  VerifiedAppActionTypes
}                                                                 from "./actions";


@Injectable()
export class AppsEffects {


  @Effect()
  loadCharts: Observable<Action> = this.actions$.pipe(
    ofType(AppCommonActionTypes.LoadCharts),
    mergeMap(
      (action: AppCommonActions) => this.http.get(`/v1/api/helm/repositories/${action.payload}/charts`),
      (action: AppCommonActions, charts) => {
        return  new LoadChartsSuccess({ repo: action.payload, charts})
      }
    ),
  );


  @Effect()
  appFilter: Observable<Action> = this.actions$.pipe(
    ofType(AppCommonActionTypes.AppFilter),
    distinctUntilChanged()
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

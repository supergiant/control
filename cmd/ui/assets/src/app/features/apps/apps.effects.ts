import { Injectable }                                     from '@angular/core';
import { Actions, Effect, ofType }                        from '@ngrx/effects';
import { distinctUntilChanged, map, mergeMap, switchMap } from 'rxjs/operators';
import { HttpClient }                                     from '@angular/common/http';
import { Action }                                         from '@ngrx/store';
import { Observable }                                     from 'rxjs';
import { Chart }                                          from './apps.reducer';
import {
  AppCommonActions,
  AppCommonActionTypes,
  AppDetailActionTypes, LoadAppDetails,
  SetAppDetails, LoadChartsSuccess,
}                                                         from './actions';


@Injectable()
export class AppsEffects {


  @Effect()
  loadCharts: Observable<Action> = this.actions$.pipe(
    ofType(AppCommonActionTypes.LoadCharts),
    mergeMap(
      (action: AppCommonActions) => this.http.get(`/v1/api/helm/repositories/${action.payload}/charts`),
      (action: AppCommonActions, charts) => {
        return new LoadChartsSuccess({ repo: action.payload, charts });
      }
    ),
  );


  @Effect()
  appFilter: Observable<Action> = this.actions$.pipe(
    ofType(AppCommonActionTypes.AppFilter),
    distinctUntilChanged()
  );

  @Effect()
  loadChartDetails: Observable<SetAppDetails> = this.actions$.pipe(
    ofType(AppDetailActionTypes.LoadAppDetails),
    map((action: LoadAppDetails) => action.payload),
    switchMap(({ repo, chart }) => this.http.get(`/v1/api/helm/repositories/${repo}/charts/${chart}`)),
    map((chart: Chart) => new SetAppDetails(chart)),
  );


  constructor(
    private actions$: Actions,
    private http: HttpClient,
  ) {
  }
}

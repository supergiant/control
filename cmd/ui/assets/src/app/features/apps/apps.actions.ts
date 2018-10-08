import { Action } from '@ngrx/store';

export enum AppStoreActionTypes {
  LoadSupergiantCharts        = '[Apps] Load SG charts',
  LoadSupergiantChartsSuccess = '[Apps] Load SG charts success'
}

export class LoadSupergiantCharts implements Action {
  readonly type = AppStoreActionTypes.LoadSupergiantCharts;
}

export class LoadSupergiantChartsSuccess implements Action {
  readonly type = AppStoreActionTypes.LoadSupergiantChartsSuccess;

  // TODO: add payload interface
  constructor(public payload: any){}
}

export type AppsActions =
  LoadSupergiantCharts |
  LoadSupergiantChartsSuccess;

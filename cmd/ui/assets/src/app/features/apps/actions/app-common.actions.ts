import { Action } from '@ngrx/store';

export enum AppCommonActionTypes {
  AppFilter         = '[AppCommon] Apps filter',
  LoadCharts        = '[AppCommon] Load charts',
  LoadChartsSuccess = '[AppCommon] Load charts success'
}

export class AppFilter implements Action {
  readonly type = AppCommonActionTypes.AppFilter;

  constructor(public payload: any) {
  }
}

export class LoadCharts implements Action {
  readonly type = AppCommonActionTypes.LoadCharts;

  // TODO: add payload interface
  constructor(public payload: any) {
  }
}

export class LoadChartsSuccess implements Action {
  readonly type = AppCommonActionTypes.LoadChartsSuccess;

  // TODO: add payload interface
  constructor(public payload: any) {
  }
}


export type AppCommonActions =
  AppFilter |
  LoadCharts |
  LoadChartsSuccess;

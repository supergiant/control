import { Action } from '@ngrx/store';

export enum AppCommonActionTypes {
  AppFilter = '[AppCommon] Apps filter'
}

export class AppFilter implements Action {
  readonly type = AppCommonActionTypes.AppFilter;

  constructor(public payload: any) {}
}

export type AppCommonActions = AppFilter;

import { Action } from '@ngrx/store';

export enum AppDetailActionTypes {
  LoadAppDetails = '[AppDetail] Load AppDetails',
  LoadAppDetailsSuccess = '[AppDetail] Load AppDetails Success'
}

export class LoadAppDetails implements Action {
  readonly type = AppDetailActionTypes.LoadAppDetails;

  constructor(public payload: any) {}

}

export class LoadAppDetailsSuccess implements Action {
  readonly type = AppDetailActionTypes.LoadAppDetailsSuccess;

  // TODO add `chartDetails` type
  constructor(public payload: any) {}
}

export type AppDetailActions =
  LoadAppDetails |
  LoadAppDetailsSuccess;

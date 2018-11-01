import { Action } from '@ngrx/store';

export enum AppDetailActionTypes {
  LoadAppDetails = '[AppDetail] Load AppDetails',
  SetAppDetails = '[AppDetail] Load AppDetails Success'
}

export class LoadAppDetails implements Action {
  readonly type = AppDetailActionTypes.LoadAppDetails;

  constructor(public payload: any) {}

}

export class SetAppDetails implements Action {
  readonly type = AppDetailActionTypes.SetAppDetails;

  // TODO add `chartDetails` type
  constructor(public payload: any) {}
}

export type AppDetailActions =
  LoadAppDetails |
  SetAppDetails;

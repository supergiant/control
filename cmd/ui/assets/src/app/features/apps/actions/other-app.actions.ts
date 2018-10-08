import { Action } from '@ngrx/store';
import { Chart }  from "../apps.reducer";

export enum OtherAppActionTypes {
  LoadOtherApps = '[OtherApp] Load OtherApps',
  LoadOtherAppsSuccess = '[OtherApp] Load OtherApps Success'
}

export class LoadOtherApps implements Action {
  readonly type = OtherAppActionTypes.LoadOtherApps;
}

export class LoadOtherAppsSuccess implements Action {
  readonly type = OtherAppActionTypes.LoadOtherAppsSuccess;

  constructor(public payload: Chart[]){}
}

export type OtherAppActions =
  LoadOtherAppsSuccess |
  LoadOtherApps;

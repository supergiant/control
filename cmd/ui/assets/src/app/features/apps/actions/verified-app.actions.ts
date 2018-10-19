import { Action } from '@ngrx/store';
import { Chart }  from "../apps.reducer";

export enum VerifiedAppActionTypes {
  LoadVerifiedApps = '[VerifiedApp] Load VerifiedApps',
  LoadVerifiedAppsSuccess = '[VerifiedApp] Load VerifiedApps success'
}

export class LoadVerifiedApps implements Action {
  readonly type = VerifiedAppActionTypes.LoadVerifiedApps;
}

export class LoadVerifiedAppsSuccess implements Action {
  readonly type = VerifiedAppActionTypes.LoadVerifiedAppsSuccess;

  // TODO: add payload interface
  constructor(public payload: Chart[]){}
}

export type VerifiedAppActions =
  LoadVerifiedAppsSuccess |
  LoadVerifiedApps;

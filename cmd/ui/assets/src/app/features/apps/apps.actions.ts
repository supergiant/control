import { Action } from '@ngrx/store';

export enum AppsActionTypes {
  LoadAppss = '[Apps] Load Appss'
}

export class LoadAppss implements Action {
  readonly type = AppsActionTypes.LoadAppss;
}

export type AppsActions = LoadAppss;

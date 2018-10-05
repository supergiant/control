import { Action } from '@ngrx/store';

export enum AppsActionTypes {
  LoadAppss = '[Apps] Load Appss',
  SetCharts = '[Apps] Set charts'
}

export class LoadAppss implements Action {
  readonly type = AppsActionTypes.LoadAppss;
}

export class SetCharts implements Action {
  readonly type = AppsActionTypes.SetCharts;

  constructor(public payload: any){}
}

export type AppsActions = LoadAppss | SetCharts;

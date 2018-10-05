import {
  ActionReducer,
  ActionReducerMap,
  createFeatureSelector,
  createSelector,
  MetaReducer
} from '@ngrx/store';
import { environment } from '../../environments/environment';
import * as fromApps from '../features/apps/apps.reducer';

export interface State {

  apps: fromApps.State;
}

export const reducers: ActionReducerMap<State> = {

  apps: fromApps.reducer,
};


export const metaReducers: MetaReducer<State>[] = !environment.production ? [] : [];

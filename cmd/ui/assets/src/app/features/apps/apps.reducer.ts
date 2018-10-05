import { Action } from '@ngrx/store';
import { AppsActions, AppsActionTypes } from './apps.actions';

export interface State {

}

export const initialState: State = {

};

export function reducer(state = initialState, action: AppsActions): State {
  switch (action.type) {

    case AppsActionTypes.LoadAppss:
      return state;


    default:
      return state;
  }
}

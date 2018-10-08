import { createSelector }                   from '@ngrx/store';
import { AppsActions, AppStoreActionTypes } from './apps.actions';
import { State }                            from '../../reducers';

export interface Chart {
  name: string;
  repo: string;
  description: string;
}

export interface AppStoreState {
  charts: {
    supergiant: Chart[]
    verified: Chart[]
  }
}

const mockCharts = [
  {
    'name': '',
    'repo': '',
    'description': '',
  },
];

export const initialState: AppStoreState = {
  charts: {
    supergiant: mockCharts,
    verified: mockCharts,
  },
};

export function reducer(state = initialState, action: AppsActions): AppStoreState {
  switch (action.type) {

    case AppStoreActionTypes.LoadSupergiantCharts:
      return state;

    case AppStoreActionTypes.LoadSupergiantChartsSuccess:
      return {
        ...state,
        charts: {
          ...state.charts,
          supergiant: action.payload
        }
      };

    default:
      return state;
  }
}


export const selectApps = createSelector(
  (state: State) => state.apps,
);
export const selectCharts = createSelector(
  selectApps,
  (state: AppStoreState) => state.charts.supergiant,
);

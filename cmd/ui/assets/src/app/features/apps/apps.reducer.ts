import { createSelector } from '@ngrx/store';
import { State }          from '../../reducers';
import {
  AppCommonActions,
  AppCommonActionTypes,
  AppDetailActions,
  AppDetailActionTypes,
}                         from "./actions";

export interface Chart {
  metadata: ChartMetadata
  values?: string;
  readme?: string;
}

export interface ChartMetadata {
  name: string;
  repo: string;
  description: string;
  version: string;
};

export interface ChartList {
  name?: string;
  description?: string;
}

export interface AppStoreState {
  charts: {
    [ key: string ]: ChartMetadata[]
  }
  currentChart: Chart
  filter: string,
}

const mockChart = {
  'name': '',
  'repo': '',
  'description': '',
  'values': ''
};


export const initialState: any = { // FIXME
  charts: {
    supergiant: [ mockChart, ],
  },
  currentChart: mockChart,
  filter: '',
};


export interface Repository {
  config: {
    url: string;
    name: string;
  }
}

// TODO: make separate reducers
type AppsActions =
  AppDetailActions |
  AppCommonActions ;

export function reducer(
  state = initialState,
  action: AppsActions
): AppStoreState {
  switch (action.type) {

    case AppDetailActionTypes.SetAppDetails:
      return {
        ...state,
        currentChart: action.payload
      };


    case AppCommonActionTypes.AppFilter:
      return {
        ...state,
        filter: action.payload
      };

    case AppCommonActionTypes.LoadChartsSuccess:
      return {
        ...state,
        charts: {
          [ action.payload.repo ]: action.payload.charts
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
  (state: AppStoreState, props: { repo }) => state.charts[ props.repo ],
);

export const selectAppDetails = createSelector(
  selectApps,
  (state: AppStoreState) => state.currentChart,
);

export const selectFilterApps = createSelector(
  selectApps,
  state => state.filter
);

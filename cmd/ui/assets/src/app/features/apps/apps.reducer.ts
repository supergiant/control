import { createSelector }                                                                     from '@ngrx/store';
import { AppStoreActionTypes, SupergiantAppActions }                                          from './actions/supergiant-app-actions';
import { State }                                                                              from '../../reducers';
import {
  AppCommonActions, AppCommonActionTypes,
  AppDetailActions,
  AppDetailActionTypes,
  VerifiedAppActions,
  VerifiedAppActionTypes
} from "./actions";
import { OtherAppActions, OtherAppActionTypes }                                               from "./actions/other-app.actions";

export interface Chart {
  name: string;
  repo: string;
  description: string;
}

export interface ChartDetails extends Chart {
  instructions?: string
}

export interface AppStoreState {
  charts: {
    supergiant: Chart[]
    verified: Chart[]
    other: Chart[]
  }
  currentChart: ChartDetails
  filter: string,
}

const mockChart = {
  'name': '',
  'repo': '',
  'description': '',
};

export const initialState: AppStoreState = {
  charts: {
    supergiant: [ mockChart, ],
    verified: [ mockChart, ],
    other: [ mockChart, ],
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
  SupergiantAppActions |
  OtherAppActions |
  AppCommonActions |
  VerifiedAppActions;

export function reducer(
  state = initialState,
  action: AppsActions
): AppStoreState {
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
      break;

    case VerifiedAppActionTypes.LoadVerifiedAppsSuccess:
      return {
        ...state,
        charts: {
          ...state.charts,
          verified: action.payload
        }
      };
      break;

    case OtherAppActionTypes.LoadOtherAppsSuccess:
      return {
        ...state,
        charts: {
          ...state.charts,
          other: action.payload
        }
      };

    case AppDetailActionTypes.LoadAppDetailsSuccess:
      return {
        ...state,
        currentChart: action.payload
      };
    case AppCommonActionTypes.AppFilter:
      return {
        ...state,
        filter: action.payload
      };

    default:
      return state;
  }
}


export const selectApps = createSelector(
  (state: State) => state.apps,
);

export const selectSupergiantCharts = createSelector(
  selectApps,
  (state: AppStoreState) => state.charts.supergiant,
);

export const selectVerifiedCharts = createSelector(
  selectApps,
  (state: AppStoreState) => state.charts.verified,
);

export const selectOtherCharts = createSelector(
  selectApps,
  (state: AppStoreState) => state.charts.other,
);

export const selectAppDetails = createSelector(
  selectApps,
  (state: AppStoreState) => state.currentChart,
);

export const selectFilterApps = createSelector(
  selectApps,
  state => state.filter
);

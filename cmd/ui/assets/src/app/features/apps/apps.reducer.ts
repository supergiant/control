import { createSelector } from '@ngrx/store';
import { AppsActions, AppsActionTypes } from './apps.actions';
import { State } from '../../reducers';

export interface Chart {
  name: string;
  repo: string;
  description: string;
}

export interface AppStoreState {
  charts: Chart[]
}

export const initialState: AppStoreState = {
  charts: [
    {
      'name': 'elasticsearch',
      'repo': 'supergiant',
      'description': 'A Helm chart for elasticsearch',
      // 'home': '',
      // 'keywords': null,
      // 'maintainers': [],
      // 'sources': null,
      // 'icon': '',
      // 'versions': [ {
      //   'version': '0.1.0',
      //   'appVersion': '',
      //   'created': '2018-02-21T11:37:37.422696-06:00',
      //   'digest': 'fba23e9ebd5c260653998f06cbd20cb85a10e0475a40d200ec33d9cebe58a962',
      //   'urls': [ 'https://supergiant.github.io/charts/elasticsearch-0.1.0.tgz' ],
      // } ],
    },
    {
      'name': 'elasticsearch',
      'repo': 'supergiant',
      'description': 'A Helm chart for elasticsearch',
    },
    {
      'name': 'elasticsearch',
      'repo': 'supergiant',
      'description': 'A Helm chart for elasticsearch',
    }

  ],

};

export function reducer(state = initialState, action: AppsActions): AppStoreState {
  switch (action.type) {

    case AppsActionTypes.LoadAppss:
      return state;

    case AppsActionTypes.SetCharts:
      return {...state, charts: action.payload};

    default:
      return state;
  }
}


export const selectApps = createSelector(
  (state: State) => state.apps,
);
export const selectCharts = createSelector(
  selectApps,
  (state: AppStoreState) => state.charts,
);

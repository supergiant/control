import { Component, OnInit }                               from '@angular/core';
import { Chart, selectFilterApps, selectSupergiantCharts } from '../../apps/apps.reducer';
import { select, Store }                                   from '@ngrx/store';
import { Observable }                                      from 'rxjs';
import { State }                                           from '../../../reducers';
import { LoadSupergiantCharts }                            from "../../apps/actions/supergiant-app-actions";
import { map, switchMap }                                  from "rxjs/operators";

@Component({
  selector: 'app-apps-supergiant',
  templateUrl: './apps-supergiant.component.html',
  styleUrls: [ './apps-supergiant.component.scss' ]
})
export class AppsSupergiantComponent implements OnInit {

  public charts$: Observable<Chart[]>;

  constructor(private store: Store<State>) {
    this.store.dispatch(new LoadSupergiantCharts());
  }

  ngOnInit() {
    this.charts$ = this.store.pipe(
      select(selectSupergiantCharts),
      switchMap((charts) => this.store.pipe(
        // filter apps if any

        select(selectFilterApps),
        map(filterMask => filterMask),
        map((filterMask) => {
          return charts.filter(chart =>
            chart.name.toLowerCase().match(filterMask) ||
            chart.description.toLocaleLowerCase().match(filterMask)
          )
        })
        )
      ),
    );
  }
}

import { Component, OnInit }             from '@angular/core';
import { ActivatedRoute }                from "@angular/router";
import { LoadAppDetails, SetAppDetails }         from "../../apps/actions";
import { State }                                 from "../../../reducers";
import { select, Store }                         from "@ngrx/store";
import { Chart, selectAppDetails }               from "../../apps/apps.reducer";
import { Observable }                            from "rxjs";
import { MatDialog }                             from "@angular/material";
import { DeployComponent }                       from "./deploy/deploy.component";
import { ConfigureComponent }                    from "./confure/configure.component";
import { map, switchMap, tap, filter, distinct } from "rxjs/operators";

@Component({
  selector: 'app-details',
  templateUrl: './app-details.component.html',
  styleUrls: ['./app-details.component.scss']
})
export class AppDetailsComponent implements OnInit {
  chartDetails$: Observable<Chart>;

  constructor(
    private route: ActivatedRoute,
    private store: Store<State>,
    public dialog: MatDialog,
  ) {

    let repo = this.route.snapshot.paramMap.get('repo');
    let chart = this.route.snapshot.paramMap.get('chart');

    this.store.dispatch(new LoadAppDetails({ repo, chart }))
  }

  ngOnInit() {
    this.chartDetails$ = this.store.pipe(select(selectAppDetails));
  }

  openDeployDialog() {
    this.dialog.open(DeployComponent, {
      data: {
        chart$: this.chartDetails$,
        routeParams: this.route.snapshot.params,
      },

    });
  }

  openConfigureDialog() {
    this.dialog.open(ConfigureComponent, {
      data: {
        chart$: this.chartDetails$,
      }
    }).afterClosed()
      .pipe(
        distinct(),
        map(
          values => this.chartDetails$.pipe(
            map(chart => {
              chart.values = values;
              return chart;
            })
          )
        ),
        switchMap(values => values),
        tap((updatedValues) => {
          this.store.dispatch(new SetAppDetails(updatedValues));
        })
      ).subscribe(val => {
      console.log('res1', val);
    });

  }
}

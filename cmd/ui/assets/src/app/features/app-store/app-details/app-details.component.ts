import { Component, OnInit }       from '@angular/core';
import { ActivatedRoute }          from "@angular/router";
import { LoadAppDetails }          from "../../apps/actions";
import { State }                   from "../../../reducers";
import { select, Store }           from "@ngrx/store";
import { Chart, selectAppDetails } from "../../apps/apps.reducer";
import { Observable }              from "rxjs";
import { MatDialog }               from "@angular/material";
import { DeployComponent }         from "./deploy/deploy.component";
import { ConfigureComponent }      from "./confure/configure.component";
import { map }                     from "rxjs/operators";

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
      data: { chart$: this.chartDetails$ }
    });
  }

  openConfigureDialog() {
    this.dialog.open(ConfigureComponent, {
      data: {
        values: this.chartDetails$.pipe(
          map(chart => chart.values)
        ),
      }
    }).afterClosed()
      .pipe(
        map(
          values => this.chartDetails$.pipe(
            map(chart => chart.values = values)
          )
        ),
        // switchMap((res) => res),
        // map(res => {
        //   return { ...res[0], values: res[1], readme: ''};
        // }),
        // map(updatedValues => this.store.dispatch(new SetAppDetails(updatedValues)))
      ).subscribe(val => {
      console.log('res1', val);
    });

  }
}

import { Component, OnInit }              from '@angular/core';
import { ActivatedRoute }                 from "@angular/router";
import { LoadAppDetails }                 from "../../apps/actions";
import { State }                          from "../../../reducers";
import { select, Store }                  from "@ngrx/store";
import { ChartDetails, selectAppDetails } from "../../apps/apps.reducer";
import { Observable }                     from "rxjs";

@Component({
  selector: 'app-details',
  templateUrl: './app-details.component.html',
  styleUrls: [ './app-details.component.scss' ]
})
export class AppDetailsComponent implements OnInit {
  chartDetails$: Observable<ChartDetails>;

  constructor(
    private route: ActivatedRoute,
    private store: Store<State>,
  ) {

    let repo = this.route.snapshot.paramMap.get('repo');
    let chart = this.route.snapshot.paramMap.get('chart');

    this.store.dispatch(new LoadAppDetails({ repo, chart }))
  }

  ngOnInit() {
    this.chartDetails$ = this.store.pipe(select(selectAppDetails));
  }

}

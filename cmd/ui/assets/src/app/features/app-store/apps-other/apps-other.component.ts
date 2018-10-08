import { Component, OnInit }                   from '@angular/core';
import { select, Store }                       from "@ngrx/store";
import { Chart, selectOtherCharts }            from "../../apps/apps.reducer";
import { State }                               from "../../../reducers";
import { Observable }                          from "rxjs";
import { LoadOtherApps } from "../../apps/actions/other-app.actions";

@Component({
  selector: 'app-apps-other',
  templateUrl: './apps-other.component.html',
  styleUrls: ['./apps-other.component.scss']
})
export class AppsOtherComponent implements OnInit {

  charts$: Observable<Chart[]>;

  constructor(private store: Store<State>) {
    this.store.dispatch(new LoadOtherApps());
  }

  ngOnInit() {
    this.charts$ = this.store.pipe(select(selectOtherCharts))
  }

}

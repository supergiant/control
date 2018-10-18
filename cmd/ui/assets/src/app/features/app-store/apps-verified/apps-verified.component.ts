import { Component, Input, OnInit }    from '@angular/core';
import { Chart, selectVerifiedCharts } from "../../apps/apps.reducer";
import { select, Store }               from "@ngrx/store";
import { State }                       from "../../../reducers";
import { Observable }                  from "rxjs";
import { LoadVerifiedApps }            from "../../apps/actions";

@Component({
  selector: 'app-apps-verified',
  templateUrl: './apps-verified.component.html',
  styleUrls: ['./apps-verified.component.scss']
})
export class AppsVerifiedComponent implements OnInit {

  public charts$: Observable<Chart[]>;

  constructor(public store: Store<State>) {
    this.store.dispatch(new LoadVerifiedApps())
  }

  ngOnInit() {
    this.charts$ = this.store.pipe(select(selectVerifiedCharts))
  }

}

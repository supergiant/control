import { Component, OnInit }             from '@angular/core';
import { Chart, selectSupergiantCharts } from '../../apps/apps.reducer';
import { select, Store }                 from '@ngrx/store';
import { Observable }                    from 'rxjs';
import { State }                         from '../../../reducers';

@Component({
  selector: 'app-apps-supergiant',
  templateUrl: './apps-supergiant.component.html',
  styleUrls: [ './apps-supergiant.component.scss' ]
})
export class AppsSupergiantComponent implements OnInit {

  public charts$: Observable<Chart[]>;

  constructor(private store: Store<State>) {
    this.charts$ = this.store.pipe(select(selectSupergiantCharts))
  }

  ngOnInit() {
    this.charts$.subscribe(charts => console.log(charts))
  }
}

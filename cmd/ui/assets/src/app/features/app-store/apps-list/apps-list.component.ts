import { Component, Input, OnInit }              from '@angular/core';
import { Chart, selectCharts, selectFilterApps } from "../../apps/apps.reducer";
import { ActivatedRoute, NavigationEnd, Router } from "@angular/router";
import { HttpClient }                            from "@angular/common/http";
import { Observable, Subscription }              from "rxjs";
import { select, Store }                         from "@ngrx/store";
import { State }                                 from "../../../reducers";
import { LoadCharts }                            from "../../apps/actions";
import { filter, map, switchMap }                from "rxjs/operators";

@Component({
  selector: 'apps-list',
  templateUrl: './apps-list.component.html',
  styleUrls: [ './apps-list.component.scss' ]
})
export class AppsListComponent implements OnInit {

  @Input() charts: Chart[] | Observable<any>;
  repo: string;
  private subscription: Subscription;

  constructor(
    private route: ActivatedRoute,
    private http: HttpClient,
    private store: Store<State>,
    public router: Router,
  ) {
    this.updateCharts();
  }

  private updateCharts() {
    this.repo = this.route.snapshot.paramMap.get('repo');

    this.store.dispatch(new LoadCharts(this.repo));
    this.charts = this.store.pipe(
      select(selectCharts, { repo: this.repo }),
      filter((charts: Chart[]) => Boolean(charts)),
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
      )
    )
  }

  ngOnInit() {
    this.subscription = this.router.events.pipe(
      filter(event => event instanceof NavigationEnd)
    ).subscribe(() => {
      this.updateCharts();
    });
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }

  removeRepo() {
    this.http.delete(`/v1/api/helm/repositories/${this.repo}`).subscribe(() => {
        // TODO: progress spinner
        this.router.navigate([ 'apps' ]);
        window.location.reload();
      }
    );
  }
}

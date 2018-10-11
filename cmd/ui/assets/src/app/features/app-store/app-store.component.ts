import { Component } from '@angular/core';
import { Router }               from "@angular/router";
import { State }                from "../../reducers";
import { select, Store }        from "@ngrx/store";
import { selectAppDetails }     from "../apps/apps.reducer";
import { Observable }           from "rxjs";

@Component({
  selector: 'app-app-store',
  templateUrl: './app-store.component.html',
  styleUrls: [ './app-store.component.scss' ]
})
export class AppStoreComponent {
  showBreadcrumbs: boolean;
  breadcrumbsData$: Observable<any>;

  constructor(
    public router: Router,
    private store: Store<State>
  ) {

    this.router.events.subscribe(() => {
      const detailPageRegexp = /\/apps\/[a-z]+\/details\/[a-z]/;
      this.showBreadcrumbs = Boolean(this.router.url.match(detailPageRegexp));
    });

    this.breadcrumbsData$ = this.store.pipe(
      select(selectAppDetails)
    );
  }
}

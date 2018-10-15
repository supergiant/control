import { Component } from '@angular/core';
import { Router }            from "@angular/router";
import { State }             from "../../reducers";
import { select, Store }     from "@ngrx/store";
import { selectAppDetails }  from "../apps/apps.reducer";
import { Observable }        from "rxjs";
import { AppFilter }         from "../apps/actions";
import { MatDialog }         from "@angular/material";
import { AppsAddComponent }  from "./apps-add/apps-add.component";

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
    private store: Store<State>,
    private dialog: MatDialog,
  ) {

    this.router.events.subscribe(() => {
      const detailPageRegexp = /\/apps\/[a-z]+\/details\/[a-z]/;
      this.showBreadcrumbs = Boolean(this.router.url.match(detailPageRegexp));
      this.store.dispatch(new AppFilter(''))
    });

    this.breadcrumbsData$ = this.store.pipe(
      select(selectAppDetails)
    );
  }

  // TODO create separate component
  filterApps(e) {
    this.store.dispatch(new AppFilter(e.target.value))
  }

  filterClear(e) {
    e.stopPropagation();
    e.stopImmediatePropagation();
    e.target.value = '';
  }

  addRepo() {
    this.dialog.open(AppsAddComponent)
  }
}

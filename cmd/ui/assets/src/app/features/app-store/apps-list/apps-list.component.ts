import { HttpClient }                                  from '@angular/common/http';
import { Component, OnInit, ViewChild }                from '@angular/core';
import { MatPaginator, MatDialog }                     from '@angular/material';
import { ActivatedRoute, NavigationEnd, Router }       from '@angular/router';
import { select, Store }                               from '@ngrx/store';
import { combineLatest, Observable, of, Subscription } from 'rxjs';
import { filter, map, startWith, switchMap, tap }      from 'rxjs/operators';
import { State }                                       from '../../../reducers';
import { LoadCharts }                                  from '../../apps/actions';
import { ChartList, selectCharts, selectFilterApps }   from '../../apps/apps.reducer';
import { RemoveRepoDialogComponent }                   from 'app/features/app-store/apps-list/remove-repo-dialog/remove-repo-dialog.component';
import { AppsAddComponent }                            from 'app/features/app-store/apps-add/apps-add.component';

@Component({
  selector: 'apps-list',
  templateUrl: './apps-list.component.html',
  styleUrls: ['./apps-list.component.scss']
})
export class AppsListComponent implements OnInit {

  paginator;
  currentPage$;
  itemsCount = 0;
  repo: string;
  charts$: Observable<any>;

  @ViewChild(MatPaginator) set matPaginator(mp: MatPaginator) {
    this.paginator = mp;
  }

  private subscription: Subscription;

  constructor(
    private route: ActivatedRoute,
    private http: HttpClient,
    private store: Store<State>,
    public router: Router,
    public dialog: MatDialog,
  ) {
    this.updateCharts();
  }

  ngOnInit() {
    this.currentPage$ = this.paginator.page.pipe(
      startWith({
        previousPageIndex: 0,
        pageIndex: 0,
        pageSize: 10,
        length: 100
      })
    );

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

    const dialogRef = this.dialog.open(RemoveRepoDialogComponent, {
      width: '300px',
    });

    const deleteRequest$ = this.http.delete(`/api/v1/helm/repositories/${this.repo}`);

    dialogRef
      .afterClosed()
      .pipe(
        filter(confirm => confirm),
        map(_ => deleteRequest$),
        switchMap(result => result)
      ).subscribe((result) => {
        // TODO: progress spinner
        // TODO: handle error
        this.router.navigate(['../'], { relativeTo: this.route });
        // window.location.reload();
      }
    );
  }

  private updateCharts() {
    this.repo = this.route.snapshot.paramMap.get('repo');

    this.store.dispatch(new LoadCharts(this.repo));

    this.charts$ = this.store.pipe(
      select(selectCharts, { repo: this.repo }),
      filter((charts: ChartList[]) => Boolean(charts)),
      switchMap((charts) => this.store.pipe(
        // filter apps if any

        select(selectFilterApps),
        map(filterMask => filterMask),
        map((filterMask) => charts.filter(chart =>
          chart.name.toLowerCase().match(filterMask) ||
          chart.description.toLocaleLowerCase().match(filterMask)
          )
        ),
        tap((charts: any[]) => this.itemsCount = charts.length),
        switchMap(charts => combineLatest(of(charts), this.currentPage$)),
        switchMap((params: any[]) => {
          const [charts, page] = params;

          const start = page.pageIndex * page.pageSize;
          const end = start + page.pageSize;

          return of(charts.slice(start, end));
        }))
        //  TODO if you wonna write some more code here then STOP

        //  breathe out...

        //  do some meditation
        //  go get your coffee
        //  and move this to *.effects.ts file
      ));
  }

  addRepo() {
    this.dialog.open(AppsAddComponent, { width: "350px" });
  }

}

import { Component, OnInit, OnDestroy, Pipe, PipeTransform } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { AppsService } from '../apps.service';

@Component({
  selector: 'app-apps-list',
  templateUrl: './apps-list.component.html',
  styleUrls: ['./apps-list.component.scss']
})
export class AppsListComponent implements OnInit, OnDestroy {
  public selected: Array<any> = [];
  public rows: Array<any> = [];
  private subscriptions = new Subscription();
  public unfilteredRows: Array<any> = [];
  public filterText: string = '';
  constructor(
    private supergiant: Supergiant,
  ) { }

  ngOnInit() {
    this.getApps();
  }

  // getCharts() {
  //   this.subscriptions.add(Observable.timer(0, 5000)
  //     .switchMap(() => this.supergiant.HelmCharts.get()).subscribe(
  //     (apps) => { this.apps = apps.items; },
  //     () => { }));
  // }

  getApps() {
    this.subscriptions.add(Observable.timer(0, 10000)
      .switchMap(() => this.supergiant.HelmReleases.get()).subscribe(
      (deployments) => {
        this.unfilteredRows = deployments.items;
        this.rows = this.filterRows(deployments.items, this.filterText);
      },
      () => { }));
  }

  onActivate(activated) {
    if (activated.type === 'click') {
      // this.router.navigate(['/apps', activated.row.id]);
    }
  }

  filterRows(filterRows: Array<any>, filterText: string): Array<any> {
    console.log(filterRows);
    console.log(filterText);
    if (filterText === '') {
      return filterRows;
    }
    const matchingRows = [];
    for (const row of filterRows) {
      for (const key of Object.keys(row)) {
        if ( row[key] != null) {
          const value = row[key].toString().toLowerCase();
          if (value.toString().indexOf(filterText.toLowerCase()) >= 0) {
            matchingRows.push(row);
            break;
          }
        }
      }
    }
    return matchingRows;
  }

  keyUpFilter(filterText) {
    this.filterText = filterText;
    this.rows = this.filterRows(this.unfilteredRows, filterText);
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }
}

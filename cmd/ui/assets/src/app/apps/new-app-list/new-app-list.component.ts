
import {timer as observableTimer,  Observable ,  Subscription } from 'rxjs';

import {switchMap} from 'rxjs/operators';
import { Component, OnInit, OnDestroy, Pipe, PipeTransform } from '@angular/core';
import { Router } from '@angular/router';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Location } from '@angular/common';

@Component({
  selector: 'app-new-app-list',
  templateUrl: './new-app-list.component.html',
  styleUrls: ['./new-app-list.component.scss']
})
export class NewAppListComponent implements OnInit, OnDestroy {
  public selected: Array<any> = [];
  public rows: Array<any> = [];
  public columns: Array<any> = [];
  private subscriptions = new Subscription();
  public unfilteredRows: Array<any> = [];
  public filterText = '';
  constructor(
    private supergiant: Supergiant,
    private router: Router,
    private location: Location,
  ) { }

  ngOnInit() {
    this.getCharts();

  }
  onActivate(activated) {
    if (activated.type === 'click') {
      this.router.navigate(['/apps/new', activated.row.id]);
    }
  }

  getCharts() {
    this.subscriptions.add(observableTimer(0, 30000).pipe(
      switchMap(() => this.supergiant.HelmCharts.get())).subscribe(
        (apps) => {
          this.unfilteredRows = apps.items;
          this.rows = this.filterRows(apps.items, this.filterText);
        },
        () => { }));
  }

  goBack() {
    this.location.back();
  }
  filterRows(filterRows: Array<any>, filterText: string): Array<any> {
    if (filterText === '') {
      return filterRows;
    }
    const matchingRows = [];
    for (const row of filterRows) {
      for (const key of Object.keys(row)) {
        if (row[key] != null) {
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
  // getApps() {
  //   this.subscriptions.add(Observable.timer(0, 10000)
  //     .switchMap(() => this.supergiant.HelmReleases.get()).subscribe(
  //     (deployments) => { this.rows = deployments.items; },
  //     () => { }));
  // }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }
}

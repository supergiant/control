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

  pApps: number[] = [];
  pDeployments: number[] = [];
  private apps = [];
  private deployments = [];
  filteredApps = [];
  filteredDeployments = [];
  subscriptions = new Subscription();
  searchString = '';

  constructor(
    private appsService: AppsService,
    private supergiant: Supergiant,
  ) { }

  ngOnInit() {
    this.searchString = this.appsService.searchString;
    this.getApps();
    this.getDeployments();
  }

  getApps() {
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.HelmCharts.get()).subscribe(
      (apps) => { this.apps = apps.items; },
      () => { }));
  }

  getDeployments() {
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.HelmReleases.get()).subscribe(
      (deployments) => { this.deployments = deployments.items; },
      () => { }));
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }
}

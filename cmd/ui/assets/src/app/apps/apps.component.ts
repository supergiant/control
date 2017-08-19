import { Component, OnInit } from '@angular/core';
import { Observable } from 'rxjs/Rx';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../shared/supergiant/supergiant.service'

@Component({
  selector: 'app-apps',
  templateUrl: './apps.component.html',
  styleUrls: ['./apps.component.css']
})
export class AppsComponent implements OnInit {
  private apps = [];
  private deployments = [];
  subscriptions = new Subscription();

  constructor(
    private supergiant: Supergiant,
  ) { }

  //get accouts when page loads
  ngOnInit() {
    this.getApps()
    this.getDeployments()
  }
  //get accounts
  getApps() {
    this.subscriptions.add(Observable.timer(0, 5000)
    .switchMap(() => this.supergiant.HelmCharts.get()).subscribe(
      (apps) => { this.apps = apps.items},
      () => {}))
  }

  getDeployments() {
    this.subscriptions.add(Observable.timer(0, 5000)
    .switchMap(() => this.supergiant.HelmReleases.get()).subscribe(
      (deployments) => { this.deployments = deployments.items},
      () => {}))
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe()
  }
}

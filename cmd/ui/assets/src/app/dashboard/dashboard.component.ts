import { Component, OnInit, OnDestroy, ViewEncapsulation } from '@angular/core';
import { of, Subscription, timer as observableTimer } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import { convertIsoToHumanReadable } from '../shared/helpers/helpers';
import { MatTableDataSource } from '@angular/material';

import { Supergiant } from '../shared/supergiant/supergiant.service';
import { AuthService } from '../shared/supergiant/auth/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class DashboardComponent implements OnInit, OnDestroy {
  public subscriptions = new Subscription();
  public cloudAccounts: any;
  public clusters: any;

  public clusterListView: string;
  public userChangedView = false;

  constructor(
    private supergiant: Supergiant,
    public auth: AuthService,
    private router: Router
  ) { }

  setClusterListView(view, e?) {
    if (e) {
      this.userChangedView = true;
      this.clusterListView = view;
    }
  }

  sortByName(cluster1, cluster2) {
    const clusterName1 = cluster1.name.toLowerCase();
    const clusterName2 = cluster2.name.toLowerCase();

    let comparison = 0;
    if (clusterName1 > clusterName2) {
      comparison = 1;
    } else if (clusterName1 < clusterName2) {
      comparison = -1;
    }
    return comparison;
  }

  getCloudAccounts() {
    this.subscriptions.add(this.supergiant.CloudAccounts.get().subscribe(
      (cloudAccounts) => {
        this.cloudAccounts = cloudAccounts;
      }));
  }

  getClusters() {
    this.subscriptions.add(observableTimer(0, 10000).pipe(
      switchMap(() => this.supergiant.Kubes.get())).subscribe(
        clusters => {
          // TODO: this is terrible
          clusters.map((c, i) => {
            
            c.index = i + 1;
            c.dataSource = new MatTableDataSource([
            {
              state: c.state,
              region: c.region,
              accountName: c.accountName,
              K8SVersion: c.K8SVersion,
              masters: Object.keys(c.masters),
              nodes: Object.keys(c.nodes),
              operatingSystem: c.operatingSystem,
              dockerVersion: c.dockerVersion,
              helmVersion: c.helmVersion,
              rbacEnabled: c.rbacEnabled,
              arch: c.arch,
            }])
          });

          this.clusters = clusters.sort(this.sortByName);

          if (!this.userChangedView) {
            if (clusters.length > 5) {
              this.clusterListView = "list"
            } else {
              this.clusterListView = "orb"
            }
          }
      },
      err => console.error(err)
    ));
  }

  trackByFn(index, item) {
    return index;
  }

  logout() {
    this.auth.logout();
    this.router.navigate(['']);
  }

  ngOnInit() {
    this.getCloudAccounts();
    this.getClusters();
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

}

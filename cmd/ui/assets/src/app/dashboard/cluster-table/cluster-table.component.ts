import { Component, OnInit, Input, OnDestroy, ChangeDetectionStrategy } from '@angular/core';
import { MatTableDataSource } from '@angular/material';
import { Router } from '@angular/router';
import { of, Subscription, timer as observableTimer } from 'rxjs';
import { switchMap } from 'rxjs/operators';

import { Supergiant } from '../../shared/supergiant/supergiant.service';

@Component({
  selector: 'cluster-table',
  templateUrl: './cluster-table.component.html',
  styleUrls: ['./cluster-table.component.scss']
})
export class ClusterTableComponent implements OnInit, OnDestroy {

  private subscriptions = new Subscription();
  public cpuUsage;
  public ramUsage;
  public clusterColumns = ["state", "accountName", "region", "cpu", "ram", "mastersCount", "nodesCount", "k8sversion", "operatingSystem", "dockerVersion", "helmVersion"];

  constructor(
    private router: Router,
    private supergiant: Supergiant,
  ) { }

  @Input() cluster: any

  updateMetrics(metrics) {
    this.cpuUsage = (metrics.cpu * 100).toFixed(1);
    this.ramUsage = (metrics.memory * 100).toFixed(1);
  }

  getMetrics() {
    this.subscriptions.add(observableTimer(0, 10000).pipe(
      switchMap(() => this.supergiant.Kubes.getClusterMetrics(this.cluster.name))).subscribe(
        res => this.updateMetrics(res),
        err => console.error(err)
      )
    )
  }

  ngOnInit() {
    this.getMetrics();
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

}

import { Component, OnInit, OnDestroy } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { MatDialogRef } from '@angular/material';
import { Subscription, timer as observableTimer } from 'rxjs';
import { switchMap } from 'rxjs/operators';

import { Supergiant } from '../../../shared/supergiant/supergiant.service';

@Component({
  selector: 'cluster-list-modal',
  templateUrl: './cluster-list-modal.component.html',
  styleUrls: ['./cluster-list-modal.component.scss'],
})
export class ClusterListModalComponent implements OnInit, OnDestroy {

  constructor(
    private dialogRef: MatDialogRef<ClusterListModalComponent>,
    private router: Router,
    private route: ActivatedRoute,
    private supergiant: Supergiant
  ) {  }

  public clusters: Array<any>;
  private subscriptions = new Subscription();

  navigate(clusterId) {
    this.router.navigate(['/clusters/', clusterId]);
    this.dialogRef.close();
  }

  newCluster() {
    this.router.navigate(['/clusters/new']);
    this.dialogRef.close();
  }

  getClusters() {
    this.subscriptions.add(observableTimer(0, 5000).pipe(
    switchMap(() => this.supergiant.Kubes.get())).subscribe(
      clusters => this.clusters = clusters,
      err => console.error(err)
    ));
  }

  ngOnInit() {
    this.getClusters();
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }
}

import { Component, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material';
import { Subscription, timer as observableTimer } from 'rxjs';
import { switchMap } from 'rxjs/operators';

import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { ClusterListModalComponent } from './cluster-list-modal/cluster-list-modal.component';

@Component({
  selector: 'app-toolbar',
  templateUrl: './toolbar.component.html',
  styleUrls: ['./toolbar.component.scss'],
})
export class ToolbarComponent implements OnInit {

  constructor(
    private supergiant: Supergiant,
    private dialog: MatDialog
  ) { }

  private subscriptions = new Subscription();
  private clusters: Array<any>;

  openClusterList(event) {
    this.initDialog(event);
  }

  initDialog(event) {
    const popupWidth = 200;
    const dialogRef = this.dialog.open(ClusterListModalComponent, {
      width: `${popupWidth}px`,
      backdropClass: 'backdrop',
      data: { clusters: this.clusters }
    });
    dialogRef.updatePosition({
      top: `${event.clientY + 20}px`
    });
    return dialogRef;
  }

  ngOnInit() {
  }
}

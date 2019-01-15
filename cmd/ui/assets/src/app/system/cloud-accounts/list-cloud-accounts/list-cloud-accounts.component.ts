import { MatDialog, MatPaginator, MatSort, MatTableDataSource } from '@angular/material';
import { Component, OnInit, OnDestroy, ViewEncapsulation, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { switchMap } from 'rxjs/operators';
import { timer as observableTimer,  Subscription } from 'rxjs';

import { Supergiant } from '../../../shared/supergiant/supergiant.service';
import { Notifications } from '../../../shared/notifications/notifications.service';
import { ConfirmModalComponent } from '../../../shared/modals/confirm-modal/confirm-modal.component';

@Component({
  selector: 'app-list-cloud-accounts',
  templateUrl: './list-cloud-accounts.component.html',
  styleUrls: ['./list-cloud-accounts.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class ListCloudAccountsComponent implements OnInit, OnDestroy {
  public subscriptions = new Subscription();
  public accounts: any;
  public accountColumns = ['provider', 'name', 'edit', 'delete'];
  public deletingAccount: string;

  constructor(
    private supergiant: Supergiant,
    private notifications: Notifications,
    private router: Router,
    private route: ActivatedRoute,
    public dialog: MatDialog,
  ) { }

  @ViewChild(MatSort) sort: MatSort;
  @ViewChild(MatPaginator) paginator: MatPaginator;

  activeCloudAccounts = new Set();

  getClusters() {
    this.subscriptions.add(observableTimer(0, 10000).pipe(
      switchMap(() => this.supergiant.Kubes.get())).subscribe(
        kubes => {
          kubes.forEach(k => this.activeCloudAccounts.add(k.accountName));
        },
        err => console.error(err)
      ));
  }

  getCloudAccounts() {
    // TODO: do we really need to poll this page?
    this.subscriptions.add(observableTimer(0, 10000).pipe(
      switchMap(() => this.supergiant.CloudAccounts.get())).subscribe(
        accounts => {
          this.accounts = new MatTableDataSource(accounts);
          this.accounts.sort = this.sort;
          this.accounts.paginator = this.paginator;
        },
        err => this.notifications.display('warn', 'Connection Issue.', err)
      ));
  }

  initDialog(event) {
    const popupWidth = 250;
    const dialogRef = this.dialog.open(ConfirmModalComponent, {
      width: `${popupWidth}px`,
    });
    dialogRef.updatePosition({
      top: `${event.clientY}px`,
      left: `${event.clientX - popupWidth - 10}px`,
    });
    return dialogRef;
  }

  delete(name, event) {
    const modal = this.initDialog(event);

    modal.afterClosed().subscribe(res => {
      if (res) {
        this.supergiant.CloudAccounts.delete(name).subscribe(
          res => {
            const refreshedAccounts = this.accounts.data.filter(account => account.name != name);
            this.accounts = new MatTableDataSource(refreshedAccounts);
            this.accounts.sort = this.sort;
            this.accounts.paginator = this.paginator;
          },
          err => this.notifications.display('warn', 'Error deleting account:', err)
        );
      }
    });
  }

  edit(name) {
    this.router.navigate(['../edit', name], { relativeTo: this.route });
  }

  ngOnInit() {
    this.getClusters();
    this.getCloudAccounts();
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

}

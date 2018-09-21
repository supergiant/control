import { Component, OnDestroy, OnInit } from '@angular/core';
import { CloudAccountsService } from '../cloud-accounts.service';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { Observable } from 'rxjs/Observable';

@Component({
  selector: 'app-cloud-accounts-list',
  templateUrl: './cloud-accounts-list.component.html',
  styleUrls: ['./cloud-accounts-list.component.scss']
})
export class CloudAccountsListComponent implements OnInit, OnDestroy {
  public p: number[] = [];
  public cloudAccounts = [];
  private subscriptions = new Subscription();
  public i: number;
  public id: number;

  constructor(
    public cloudAccountsService: CloudAccountsService,
    private supergiant: Supergiant,
    private notifications: Notifications,
  ) { }


  ngOnInit() {
    this.getAccounts();
  }

  getAccounts() {
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.CloudAccounts.get()).subscribe(
      (cloudAccounts) => { this.cloudAccounts = cloudAccounts.items; },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }
}

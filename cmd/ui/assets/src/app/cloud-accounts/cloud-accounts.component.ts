import { Component, OnDestroy, OnInit } from '@angular/core';
import { CloudAccountsService } from './cloud-accounts.service';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../shared/supergiant/supergiant.service'
import { Notifications } from '../shared/notifications/notifications.service'
import { Observable } from 'rxjs/Rx';


@Component({
  selector: 'app-cloud-accounts',
  templateUrl: './cloud-accounts.component.html',
  styleUrls: ['./cloud-accounts.component.css']
})
export class CloudAccountsComponent implements OnInit {
  cloudAccounts: any;
  private subscription: Subscription;

  constructor(
    private cloudAccountsService: CloudAccountsService,
    private supergiant: Supergiant,
    private notifications: Notifications,
  ) { }

  //get accouts when page loads
  ngOnInit() {
    this.getAccounts()
  }
  //get accounts
  getAccounts() {
    this.subscription = Observable.timer(0, 10000)
    .switchMap(() => this.supergiant.CloudAccounts.get()).subscribe(
      cloudAccount=>{ this.cloudAccounts = cloudAccount.json().items},
      (err) =>{ this.notifications.display("warn", "Connection Issue.", err)})
  }

  ngOnDestroy(){
    this.subscription.unsubscribe();
  }
}

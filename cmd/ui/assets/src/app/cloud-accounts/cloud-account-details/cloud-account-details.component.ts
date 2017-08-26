import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { ActivatedRoute, Router } from '@angular/router';
import { Observable } from 'rxjs/Observable';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';

@Component({
  selector: 'app-cloud-account-details',
  templateUrl: './cloud-account-details.component.html',
  styleUrls: ['./cloud-account-details.component.css']
})
export class CloudAccountDetailsComponent implements OnInit, OnDestroy {
  id: number;
  subscriptions = new Subscription();
  cloudAccount: any;
  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private supergiant: Supergiant,
    private notifications: Notifications,
  ) { }

  ngOnInit() {
    this.id = this.route.snapshot.params.id;
    this.getAccount();
  }

  getAccount() {
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.CloudAccounts.get(this.id)).subscribe(
      (cloudAccount) => { this.cloudAccount = cloudAccount; },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  goBack() {
    this.router.navigate(['/cloud-accounts']);
  }
  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

}

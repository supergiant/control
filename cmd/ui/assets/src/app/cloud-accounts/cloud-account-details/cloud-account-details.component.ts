
import {timer as observableTimer,  Subscription ,  Observable } from 'rxjs';

import {switchMap} from 'rxjs/operators';
import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { SystemModalService } from '../../shared/system-modal/system-modal.service';
import { LoginComponent } from '../../login/login.component';


@Component({
  selector: 'app-cloud-account-details',
  templateUrl: './cloud-account-details.component.html',
  styleUrls: ['./cloud-account-details.component.scss']
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
    private systemModalService: SystemModalService,
    public loginComponent: LoginComponent,
  ) { }

  ngOnInit() {
    this.id = this.route.snapshot.params.id;
    this.getAccount();
  }

  openSystemModal(message) {
    this.systemModalService.openSystemModal(message);
  }

  getAccount() {
    this.subscriptions.add(observableTimer(0, 5000).pipe(
      switchMap(() => this.supergiant.CloudAccounts.get(this.id))).subscribe(
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

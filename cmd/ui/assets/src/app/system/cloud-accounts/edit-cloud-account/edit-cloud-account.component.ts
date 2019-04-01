import { Component, OnInit, OnDestroy, ViewEncapsulation } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription } from 'rxjs';

import { Supergiant } from '../../../shared/supergiant/supergiant.service';
import { Notifications } from '../../../shared/notifications/notifications.service';

@Component({
  selector: 'app-edit-cloud-account',
  templateUrl: './edit-cloud-account.component.html',
  styleUrls: ['./edit-cloud-account.component.scss'],
  encapsulation: ViewEncapsulation.None
})

export class EditCloudAccountComponent implements OnInit, OnDestroy {

  subscriptions = new Subscription();
  public accountName: string;
  public account: any;
  public saving = false;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private supergiant: Supergiant,
    private notifications: Notifications
  ) { }

  getAccount(name) {
    this.subscriptions.add(this.supergiant.CloudAccounts.get(name).subscribe(
      data => this.account = data,
      err => console.error(err)
    ));
  }

  editAccount(account) {
    if (!this.saving) {
      this.saving = true;

      this.supergiant.CloudAccounts.update(this.accountName, account).subscribe(
        res => {
          this.displaySuccess(this.accountName);
          this.router.navigate(['../../'], { relativeTo: this.route});
        },
        err => this.displayError(this.accountName, err)
      );
    }
  }

  success(account) {
    this.notifications.display(
      'success',
      'Account: ' + account.name,
      'Updated',
    );
  }

  error(account, err) {
    console.log(err.error.devMessage);
    this.notifications.display(
      'error',
      'Account: ' + account.name,
      'Error: ' + err.error.userMessage
    );
  }

  cancel() {
    this.router.navigate(['../../'], { relativeTo: this.route});
  }

  displaySuccess(accountName) {
    this.notifications.display(
      "success",
      "Account: " + accountName,
      "Saved successfully"
    )
  }

  displayError(accountName, err) {
    let msg: string;

    if (err.error.userMessage) {
      console.log(err.error.devMessage);
      msg = err.error.userMessage
    } else {
      msg = err.error
    }

    this.notifications.display(
      'error',
      'Account: ' + accountName,
      'Error: ' + msg
    );
  }

  ngOnInit() {
    this.accountName = this.route.snapshot.params.id;
    this.getAccount(this.accountName);
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

}

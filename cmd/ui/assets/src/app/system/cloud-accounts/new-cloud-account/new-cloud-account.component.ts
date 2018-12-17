import { Component, OnInit, ViewEncapsulation, OnDestroy } from '@angular/core';
import { CloudAccountModel } from '../cloud-accounts.model';
import { Subscription } from 'rxjs';
import { Supergiant } from '../../../shared/supergiant/supergiant.service';
import { Notifications } from '../../../shared/notifications/notifications.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-new-cloud-account',
  templateUrl: './new-cloud-account.component.html',
  styleUrls: ['./new-cloud-account.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class NewCloudAccountComponent implements OnInit, OnDestroy {
  private providersObj = new CloudAccountModel;
  private subscriptions = new Subscription();
  public providers = this.providersObj.providers;
  private model: any;
  public schema: any;
  public nameIsBlank: boolean;
  public gceServiceAccountKeyIsBlank: boolean;

  public selectedProvider: string;
  private cloudAccountName: string;

  constructor(
    private supergiant: Supergiant,
    private notifications: Notifications,
    private router: Router,
  ) { }

  ngOnInit() {
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  back() {
    this.model = null;
    this.schema = null;
  }

  checkForBlankName(name) {
    console.log('name: ', name);
    console.log('type: ', typeof(name) === 'string');
    if (name) {
      console.log('true');
      this.nameIsBlank = false;
    } else {
      console.log('false');
      this.nameIsBlank = true;
    }
  }

  create(model) {
    // TODO: find a better way to do this...
    if (model.provider === 'gce') {
      const serviceAccountKey = model.credentials.service_account_key;

      if (serviceAccountKey == '') {
        this.gceServiceAccountKeyIsBlank = true;
      } else {
        this.gceServiceAccountKeyIsBlank = false;
        model.credentials = JSON.parse(model.credentials.service_account_key);
      }
    }

    this.subscriptions.add(this.supergiant.CloudAccounts.create(model).subscribe(
      data => {
        this.success(model);
        this.router.navigate(['/clusters/new']);
      },
      err => this.error(model, err)
    ));
  }

  success(model) {
    this.notifications.display(
      'success',
      'Account: ' + model.name,
      'Created',
    );
  }

  error(model, err) {
    console.log(err.error.devMessage);
    this.notifications.display(
      'error',
      'Account: ' + model.name,
      'Error: ' + err.error.userMessage
    );
  }
}

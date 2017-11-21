import { Component, OnInit, ViewEncapsulation, OnDestroy } from '@angular/core';
import { CloudAccountModel } from '../cloud-accounts.model';
import { Subscription } from 'rxjs/Subscription';
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
  private providers = [];
  private model: any;
  public schema: any;

  constructor(
    private supergiant: Supergiant,
    private notifications: Notifications,
    private router: Router,
  ) { }

  ngOnInit() {
    this.get();
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  back() {
    this.model = null;
    this.schema = null;
  }

  get() {
    for (const key in this.providersObj.providers) {
      if (key) {
        this.providers.push(key);
      }
    }
  }

  create(model) {
    this.subscriptions.add(this.supergiant.CloudAccounts.create(model).subscribe(
      (data) => {
        this.success(model);
        this.router.navigate(['/system/cloud-accounts']);
      },
      (err) => { this.error(model, err); }));
  }

  success(model) {
    this.notifications.display(
      'success',
      'Kube: ' + model.name,
      'Created...',
    );
  }

  error(model, data) {
    this.notifications.display(
      'error',
      'Kube: ' + model.name,
      'Error:' + data.statusText);
  }

  sendChoice(choice) {
    console.log(choice);
    switch (choice) {
      case 'AWS - Amazon Web Services': {
        this.model = this.providersObj.aws.model;
        this.schema = this.providersObj.aws.schema;
        break;
      }
      case 'Digital Ocean': {
        this.model = this.providersObj.digitalocean.model;
        this.schema = this.providersObj.digitalocean.schema;
        break;
      }
      case 'Packet.net': {
        this.model = this.providersObj.packet.model;
        this.schema = this.providersObj.packet.schema;
        break;
      }
      case 'OpenStack': {
        this.model = this.providersObj.openstack.model;
        this.schema = this.providersObj.openstack.schema;
        break;
      }
      case 'GCE - Google Compute Engine': {
        this.model = this.providersObj.gce.model;
        this.schema = this.providersObj.gce.schema;
        break;
      }
      default: {
        this.model = null;
        this.schema = null;
        break;
      }
    }


  }

}

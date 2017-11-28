import { Component, OnDestroy, AfterViewInit, ViewEncapsulation } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { ClusterAWSModel } from '../cluster.aws.model';
import { ClusterDigitalOceanModel } from '../cluster.digitalocean.model';
import { ClusterGCEModel } from '../cluster.gce.model';
import { ClusterOpenStackModel } from '../cluster.openstack.model';
import { ClusterPacketModel } from '../cluster.packet.model';
import { Notifications } from '../../shared/notifications/notifications.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-new-cluster',
  templateUrl: './new-cluster.component.html',
  styleUrls: ['./new-cluster.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class NewClusterComponent implements OnDestroy, AfterViewInit {
  subscriptions = new Subscription();
  cloudAccountsList = [];
  providers = [];
  awsModel = new ClusterAWSModel;
  doModel = new ClusterDigitalOceanModel;
  gceModel = new ClusterGCEModel;
  osModel = new ClusterOpenStackModel;
  packModel = new ClusterPacketModel;
  model: any;
  schema: any;


  constructor(
    private supergiant: Supergiant,
    private notifications: Notifications,
    private router: Router,
  ) { }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }


  ngAfterViewInit() {
    this.subscriptions.add(this.supergiant.CloudAccounts.get().subscribe(
      (data) => { this.providers = data.items; }
    ));
  }

  back() {
    this.model = null;
    this.schema = null;
  }

  createKube(model) {
    this.subscriptions.add(this.supergiant.Kubes.create(model).subscribe(
      (data) => {
        this.success(model);
        this.router.navigate(['/clusters']);
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
    switch (choice.provider) {
      case 'aws': {
        this.model = this.awsModel.aws.model;
        this.schema = this.awsModel.aws.schema;
        this.model.cloud_account_name = choice.name;
        break;
      }
      case 'digitalocean': {
        this.model = this.doModel.digitalocean.model;
        this.schema = this.doModel.digitalocean.schema;
        this.model.cloud_account_name = choice.name;
        break;
      }
      case 'packet': {
        this.model = this.packModel.packet.model;
        this.schema = this.packModel.packet.schema;
        this.model.cloud_account_name = choice.name;
        break;
      }
      case 'openstack': {
        this.model = this.osModel.openstack.model;
        this.schema = this.osModel.openstack.schema;
        this.model.cloud_account_name = choice.name;
        break;
      }
      case 'gce': {
        this.model = this.gceModel.gce.model;
        this.schema = this.gceModel.gce.schema;
        this.model.cloud_account_name = choice.name;
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

import { Component, OnInit, OnDestroy, AfterViewInit, ViewEncapsulation } from '@angular/core';
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
export class NewClusterComponent implements OnInit, OnDestroy, AfterViewInit {
  subscriptions = new Subscription();
  cloudAccountsList: any;
  awsModel = new ClusterAWSModel;
  doModel = new ClusterDigitalOceanModel;
  gceModel = new ClusterGCEModel;
  osModel = new ClusterOpenStackModel;
  packModel = new ClusterPacketModel;
  hasCluster = false;
  hasCloudAccount = false;
  hasApp = false;
  appCount = 0;
  data: any;
  schema: any;
  layout: any;

  clusterName: string;
  kubeProfiles = ["1234", "7878", "1092", "4555"];
  selectedKubeProfile: any;
  cloudAccounts: Array<any>;
  selectedCloudAccount: any;

  constructor(
    private supergiant: Supergiant,
    private notifications: Notifications,
    private router: Router,
  ) { }

  getCloudAccounts() {
    this.subscriptions.add(this.supergiant.CloudAccounts.get().subscribe(
      (cloudAccounts) => {
        this.cloudAccounts = cloudAccounts;
      })
    );
  }

  back() {
    this.data = null;
    this.schema = null;
  }

  createCluster(name, account, profile) {
    this.subscriptions.add(this.supergiant.Kubes.create({"clusterName": name, "profileId": profile, "cloudAccountName": account})
      .subscribe(
        // temporary objs sent to notifications
        (res) => {
          console.log(res);
          this.success({"name": name});
          this.router.navigate(['/dashboard']);
        },
        (err) => {
          console.log(err);
          this.error({"name": name}, err);
        }))
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
      'Kube: ' + this.data.name,
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
    switch (choice.provider) {
      case 'aws': {
        this.data = this.awsModel.aws.data;
        this.schema = this.awsModel.aws.schema;
        this.layout = this.awsModel.aws.layout;
        this.data.cloud_account_name = choice.name;
        break;
      }
      case 'digitalocean': {
        this.data = this.doModel.digitalocean.data;
        this.schema = this.doModel.digitalocean.schema;
        this.layout = this.doModel.digitalocean.layout;
        this.data.cloud_account_name = choice.name;
        break;
      }
      case 'packet': {
        this.data = this.packModel.packet.data;
        this.schema = this.packModel.packet.schema;
        this.layout = this.packModel.packet.layout;
        this.data.cloud_account_name = choice.name;
        break;
      }
      case 'openstack': {
        this.data = this.osModel.openstack.data;
        this.schema = this.osModel.openstack.schema;
        this.layout = this.osModel.openstack.layout;
        this.data.cloud_account_name = choice.name;
        break;
      }
      case 'gce': {
        this.data = this.gceModel.gce.data;
        this.schema = this.gceModel.gce.schema;
        this.layout = this.gceModel.gce.layout;
        this.data.cloud_account_name = choice.name;
        break;
      }
      default: {
        this.data = null;
        this.schema = null;
        this.layout = null;
        break;
      }
    }


  }

  ngOnInit() {
    this.getCloudAccounts();
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  ngAfterViewInit() {
    this.subscriptions.add(this.supergiant.CloudAccounts.get().subscribe(
      (data) => { this.cloudAccountsList = data; }
    ));
  }

}

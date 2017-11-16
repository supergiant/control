import { Component, OnDestroy, AfterViewInit, ViewEncapsulation } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { ClusterAWSModel } from '../cluster.aws.model';
import { ClusterDigitalOceanModel } from '../cluster.digitalocean.model';
import { ClusterGCEModel } from '../cluster.gce.model';
import { ClusterOpenStackModel } from '../cluster.openstack.model';
import { ClusterPacketModel } from '../cluster.packet.model';
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

  sendChoice(choice) {
    console.log(choice);
    switch (choice.provider) {
      case 'aws': {
        this.model = this.awsModel.aws.model;
        this.schema = this.awsModel.aws.schema;
        break;
      }
      case 'digitalocean': {
        this.model = this.doModel.digitalocean.model;
        this.schema = this.doModel.digitalocean.schema;
        break;
      }
      case 'packet': {
        this.model = this.packModel.packet.model;
        this.schema = this.packModel.packet.schema;
        break;
      }
      case 'openstack': {
        this.model = this.osModel.openstack.model;
        this.schema = this.osModel.openstack.schema;
        break;
      }
      case 'gce': {
        this.model = this.gceModel.gce.model;
        this.schema = this.gceModel.gce.schema;
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

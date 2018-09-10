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
export class NewClusterComponent implements OnInit, OnDestroy {
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


  // temp for demo
  clusterName: string;
  availableCloudAccounts: Array<any>;
  selectedCloudAccount: any;
  availableRegions: Array<any>;
  selectedRegion: any;
  availableMachineTypes: Array<any>;
  machineSizes: any;

  machines = [{
    machineType: null,
    role: "Master",
    qty: 1
  }];

  blankMachine = {
    machineType: null,
    role: null,
    qty: 1
  };

  profileOptions = {
    archs: ["amd64"],
    flannelVersions: ["0.10.0"],
    operatingSystems: ["linux"],
    networkTypes: ["vxlan"],
    ubuntuVersions: ["xenial"],
    helmVersions: ["2.8.0"],
    dockerVersions: ["17.06.0"],
    K8SVersions: ["1.11.1"],
    rbacEnabled: [true, false]
  }

  newDigitalOceanCluster = {
    profile: {
      masterProfiles: [],
      nodesProfiles: [],
      provider: "digitalocean",
      // will have to set this on submit for now UGH
      // region: this.selectedRegion.id,
      arch: "amd64",
      operatingSystem: "linux",
      ubuntuVersion: "xenial",
      dockerVersion: "17.06.0",
      K8SVersion: "1.11.1",
      flannelVersion: "0.10.0",
      networkType: "vxlan",
      cidr: "10.0.0.0/24",
      helmVersion: "2.8.0",
      rbacEnabled: false
    }
  }

  constructor(
    private supergiant: Supergiant,
    private notifications: Notifications,
    private router: Router,
  ) { }

  getCloudAccounts() {
    this.subscriptions.add(this.supergiant.CloudAccounts.get().subscribe(
      (cloudAccounts) => {
        this.availableCloudAccounts = cloudAccounts;
      })
    );
  }

  back() {
    this.data = null;
    this.schema = null;
  }

  compileProfiles(machines, role) {
    const filteredMachines = machines.filter(m => m.role == role);
    const compiledProfiles = [];

    filteredMachines.forEach(m => {
      for (var i = 0; i < m.qty; i++) {
        compiledProfiles.push({ image: "ubuntu-16-04-x64", size: m.machineType })
      }
    })

    return compiledProfiles;
  }

  createCluster(model) {
    // temp for demo
    model.cloudAccountName = this.selectedCloudAccount.name;
    model.clusterName = this.clusterName;
    model.profile.region = this.selectedRegion.id
    model.profile.masterProfiles = this.compileProfiles(this.machines, "Master");
    model.profile.nodesProfiles = this.compileProfiles(this.machines, "Node");

    console.log(model);

    this.subscriptions.add(this.supergiant.Kubes.create(model).subscribe(
      (data) => {
        this.success(model);
        this.router.navigate(['/clusters/', this.clusterName]);
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
    console.log("model:", model);
    console.log("data:", data);
    this.notifications.display(
      'error',
      'Kube: ' + model.name,
      'Error:' + data.statusText);
  }

  selectRegion(region) {
    this.availableMachineTypes = region.AvailableSizes;
    if (this.machines.length === 0) {
      this.machines.push(this.blankMachine);
    }
  }

  addBlankMachine() {
    this.machines.push(this.blankMachine);
  }

  deleteMachine(idx) {
    console.log(idx);
    this.machines.splice(idx, 1);
  }

  selectCloudAccount(cloudAccount) {
    switch (cloudAccount.provider) {
      case 'aws': {
        this.data = this.awsModel.aws.data;
        this.schema = this.awsModel.aws.schema;
        this.layout = this.awsModel.aws.layout;
        this.data.cloudAccountName = cloudAccount.name;
        break;
      }
      case 'digitalocean': {
        this.schema = this.doModel.digitalocean.schema;
        this.data = this.doModel.digitalocean.data;
        // this.layout = this.doModel.digitalocean.layout;
        this.data.cloudAccountName = cloudAccount.name;
        break;
      }
      case 'packet': {
        this.data = this.packModel.packet.data;
        this.schema = this.packModel.packet.schema;
        this.layout = this.packModel.packet.layout;
        this.data.cloudAccountName = cloudAccount.name;
        break;
      }
      case 'openstack': {
        this.data = this.osModel.openstack.data;
        this.schema = this.osModel.openstack.schema;
        this.layout = this.osModel.openstack.layout;
        this.data.cloudAccountName = cloudAccount.name;
        break;
      }
      case 'gce': {
        this.data = this.gceModel.gce.data;
        this.schema = this.gceModel.gce.schema;
        this.layout = this.gceModel.gce.layout;
        this.data.cloudAccountName = cloudAccount.name;
        break;
      }
      default: {
        this.data = null;
        this.schema = null;
        this.layout = null;
        break;
      }
    };

    this.subscriptions.add(this.supergiant.CloudAccounts.getRegions(cloudAccount.name).subscribe(
        regionList => {
          this.availableRegions = regionList;
          this.machineSizes = regionList.sizes;
        },
        err => this.error({}, err)
    ))
  }

  ngOnInit() {
    this.getCloudAccounts();
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

}

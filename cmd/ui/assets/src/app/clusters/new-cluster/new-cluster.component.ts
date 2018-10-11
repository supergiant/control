// EXAMPLE APP STATE

// {
//   cloudAccount: {},
//   clusterConfig: {
//     name: "...",
//     k8sVersion: "...",
//     flannelVersion: "...",
//     helmVersion: "...",
//     dockerVersion: "...",
//     ubuntuVersion: "...",
//     networkType: "...",
//     cidr: "...",
//     operatingSystem: "...",
//     arch: "..."
//   },
//   providerConfig: {
//     digitalOcean: {
//       region: "..."
//     },
//     aws: {
//       region: "...",
//       vpc: "...",
//       mastersSubnet: "...",
//       nodesSubnet: "...",
//       mastersSecurityGroup: "...",
//       nodessSecurityGroup: "..."
//     }
//   },
//   machinesConfig: [
//     {
//       role: "...",
//       type: "...",
//       qty: "..."
//     },
//     {
//       role: "...",
//       type: "...",
//       qty: "..."
//     }
//   ]
// }

import { Component, OnInit, OnDestroy, AfterViewInit, ViewEncapsulation } from '@angular/core';
import { Subscription } from 'rxjs';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { Router } from '@angular/router';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';

@Component({
  selector: 'app-new-cluster',
  templateUrl: './new-cluster.component.html',
  styleUrls: ['./new-cluster.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class NewClusterComponent implements OnInit, OnDestroy {
  subscriptions = new Subscription();

  clusterName: string;
  availableCloudAccounts: Array<any>;
  selectedCloudAccount: any;
  availableRegions: Array<any>;
  selectedRegion: any;
  availableMachineTypes: Array<any>;

  machines = [{
    machineType: null,
    role: "Master",
    qty: 1
  }];

  clusterOptions = {
    archs: ["amd64"],
    flannelVersions: ["0.10.0"],
    operatingSystems: ["linux"],
    networkTypes: ["vxlan"],
    ubuntuVersions: ["xenial"],
    helmVersions: ["2.8.0"],
    dockerVersions: ["17.06.0"],
    K8sVersions: ["1.11.1"],
    rbacEnabled: [false]
  };

  provisioning = false;
  clusterConfig: FormGroup;
  providerConfig: FormGroup;

  constructor(
    private supergiant: Supergiant,
    private notifications: Notifications,
    private router: Router,
    private formBuilder: FormBuilder
  ) { }


  getCloudAccounts() {
    this.subscriptions.add(this.supergiant.CloudAccounts.get().subscribe(
      (cloudAccounts) => {
        this.availableCloudAccounts = cloudAccounts;
      })
    );
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

  createCluster() {
    // compile frontend new-cluster model into api format
    const newClusterData:any = {};
    newClusterData.profile = this.clusterConfig.value;

    newClusterData.cloudAccountName = this.selectedCloudAccount.name;
    newClusterData.clusterName = this.clusterName;
    newClusterData.profile.region = this.providerConfig.value.region.id;
    newClusterData.profile.provider = this.selectedCloudAccount.provider;
    newClusterData.profile.masterProfiles = this.compileProfiles(this.machines, "Master");
    newClusterData.profile.nodesProfiles = this.compileProfiles(this.machines, "Node");

    this.provisioning = true;
    this.subscriptions.add(this.supergiant.Kubes.create(newClusterData).subscribe(
      (data) => {
        this.success(newClusterData);
        this.router.navigate(['/clusters/', newClusterData.clusterName]);
        this.provisioning = false;
      },
      (err) => {
        this.error(newClusterData, err);
        this.provisioning = false;
      }
    ));
  }

  success(model) {
    this.notifications.display(
      'success',
      'Kube: ' + model.clusterName,
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
      this.machines.push({
        machineType: null,
        role: null,
        qty: 1
      });
    }
  }

  addBlankMachine() {
    this.machines.push({
      machineType: null,
      role: null,
      qty: 1
    })
  }

  deleteMachine(idx) {
    this.machines.splice(idx, 1);
  }

  selectCloudAccount(cloudAccount) {
    this.selectedCloudAccount = cloudAccount

    switch (this.selectedCloudAccount.provider) {
      case "digitalocean": {
        this.providerConfig = this.formBuilder.group({
          region: [""]
        });
      }
    }

    this.subscriptions.add(this.supergiant.CloudAccounts.getRegions(cloudAccount.name).subscribe(
        regionList => {
          this.availableRegions = regionList;
        },
        err => this.error({}, err)
    ))
  }

  ngOnInit() {
    // build cluster config w/ defaults
    this.clusterConfig = this.formBuilder.group({
      K8sVersion: ["1.11.1"],
      flannelVersion: ["0.10.0"],
      helmVersion: ["2.8.0"],
      dockerVersion: ["17.06.0"],
      ubuntuVersion: ["xenial"],
      networkType: ["vxlan"],
      cidr: ["10.0.0.0/24"],
      operatingSystem: ["linux"],
      arch: ["amd64"],
      rbacEnabled: [false]
    });

    // get cloud accounts
    this.getCloudAccounts();
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

}

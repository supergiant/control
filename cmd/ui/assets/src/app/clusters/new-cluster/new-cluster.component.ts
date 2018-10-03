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
  cloudAccountsList: any;
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

  clusterOptions = {
    archs: ["amd64"],
    flannelVersions: ["0.10.0"],
    operatingSystems: ["linux"],
    networkTypes: ["vxlan"],
    ubuntuVersions: ["xenial"],
    helmVersions: ["2.8.0"],
    dockerVersions: ["17.06.0"],
    k8sVersions: ["1.11.1", "1.8"],
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
    private formBuilder: FormBuilder
  ) { }

  isLinear = false;
  clusterConfig: FormGroup;
  providerConfig: FormGroup;
  machinesConfig: FormGroup;

  getCloudAccounts() {
    this.subscriptions.add(this.supergiant.CloudAccounts.get().subscribe(
      (cloudAccounts) => {
        this.availableCloudAccounts = cloudAccounts;
        // set this.providerConfig here
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

    // this.subscriptions.add(this.supergiant.Kubes.create(model).subscribe(
    //   (data) => {
    //     this.success(model);
    //     this.router.navigate(['/clusters/', this.clusterName]);
    //   },
    //   (err) => { this.error(model, err); }));
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

    this.subscriptions.add(this.supergiant.CloudAccounts.getRegions(cloudAccount.name).subscribe(
        regionList => {
          this.availableRegions = regionList;
          this.machineSizes = regionList.sizes;
        },
        err => this.error({}, err)
    ))
  }

  ngOnInit() {
    // build 4 form groups
    // leave providerConfig empty
    this.clusterConfig = this.formBuilder.group({
      name: [""],
      k8sVersion: [""],
      flannelVersion: [""],
      helmVersion: [""],
      dockerVersion: [""],
      ubuntuVersion: [""],
      networkType: [""],
      cidr: [""],
      operatingSystem: [""],
      arch: [""]
    });

    // get cloud accounts
    this.getCloudAccounts();
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

}

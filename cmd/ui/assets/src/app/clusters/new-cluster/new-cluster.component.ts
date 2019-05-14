import { Component, OnDestroy, OnInit, ViewEncapsulation, ViewChild } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ValidatorFn, AbstractControl } from '@angular/forms';
import { Router } from '@angular/router';
import { MatHorizontalStepper, MatOption, MatSelect, MatDialog } from '@angular/material';
import { Subscription, Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { Notifications } from '../../shared/notifications/notifications.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { NodeProfileService } from '../node-profile.service';
import { CLUSTER_OPTIONS } from './cluster-options.config';
import {
  DEFAULT_MACHINE_SET,
  BLANK_MACHINE_TEMPLATE,
} from 'app/clusters/new-cluster/new-cluster.component.config';
import { sortDigitalOceanMachineTypes } from 'app/clusters/new-cluster/new-cluster.helpers';
import { IMachineType } from './new-cluster.component.interface';
import { PublicKeyModalComponent } from './public-key-modal/public-key-modal.component';

// compiler hack
declare var require: any;
const cidrRegex = require('cidr-regex');

enum MachineRoles {
  master = 'Master',
  node = 'Node'
}

enum CloudProviders {
  aws = 'aws',
  digitalocean = 'digitalocean',
  gce = 'gce',
  azure = 'azure',
}

export enum StepIndexes {
  ClusterConfig = 0,
  ProvideConfig = 1,
  MachinesConfig = 2,
  Review = 3,
}

@Component({
  selector: 'app-new-cluster',
  templateUrl: './new-cluster.component.html',
  styleUrls: ['./new-cluster.component.scss'],
  encapsulation: ViewEncapsulation.None,
})
export class NewClusterComponent implements OnInit, OnDestroy {
  subscriptions = new Subscription();

  clusterName: string;
  availableCloudAccounts$: Observable<any[]>;
  selectedCloudAccount: any;
  availableRegions: any;
  availableRegionNames: Array<string>;
  availableMachineTypes: Array<any>;
  regionsLoading = false;
  machinesLoading = false;

  // aws vars
  availabilityZones: Array<any>;
  azsLoading = false;

  machines = DEFAULT_MACHINE_SET;

  clusterOptions = CLUSTER_OPTIONS;

  machinesConfigValid: boolean;
  displayMachinesConfigWarning: boolean;
  provisioning = false;
  clusterConfig: FormGroup;
  providerConfig: FormGroup;
  unavailableClusterNames = new Set();
  regionsFilter = '';

  exposedAddressesArray: Array<string> = new Array<string>();

  @ViewChild(MatHorizontalStepper) stepper: MatHorizontalStepper;
  @ViewChild('selectedMachineType') selectedMachineType: MatSelect;

  constructor(
    private supergiant: Supergiant,
    private notifications: Notifications,
    private router: Router,
    private formBuilder: FormBuilder,
    private nodesService: NodeProfileService,
    public dialog: MatDialog,
  ) {
  }

  ngOnInit() {
    this.getClusters();
    this.availableCloudAccounts$ = this.supergiant.CloudAccounts.get().pipe(
      map(cloudAccounts => cloudAccounts.sort()),
    );

    this.clusterConfig = this.formBuilder.group({
      name: ['', [
        Validators.required,
        this.uniqueClusterName(this.unavailableClusterNames),
        Validators.maxLength(12),
        Validators.pattern('^[A-Za-z]([-A-Za-z0-9\-]*[A-Za-z0-9\-])?$')]],
      cloudAccount: ['', Validators.required],
      K8sVersion: ['1.14.0', Validators.required],
      networkProvider: ['Flannel', Validators.required],
      helmVersion: ['2.11.0', Validators.required],
      dockerVersion: ['18.06.3', Validators.required],
      ubuntuVersion: ['xenial', Validators.required],
      networkType: ['vxlan', Validators.required],
      cidr: ['10.100.0.0/16', [Validators.required, this.validCidr()]],
      operatingSystem: ['linux', Validators.required],
      arch: ['amd64', Validators.required],
      exposedAddresses: ['', this.allValidCidrs()],
    });
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  get name() {
    return this.clusterConfig.get('name');
  }

  get cidr() {
    return this.clusterConfig.get('cidr');
  }

  get exposedAddresses() {
    return this.clusterConfig.get('exposedAddresses');
  }

  get vpcCidr() {
    if (this.selectedCloudAccount && this.selectedCloudAccount.provider == 'aws') {
      return this.providerConfig.get('vpcCidr');
    } else {
      return true;
    }
  }

  get azureVNetCIDR() {
    if (this.selectedCloudAccount && this.selectedCloudAccount.provider == 'azure') {
      return this.providerConfig.get('azureVNetCIDR');
    } else {
      return true;
    }
  }

  getClusters() {
    this.supergiant.Kubes.get().subscribe(
      clusters => clusters.map(c => this.unavailableClusterNames.add(c.name)),
      err => console.error(err),
    );
  }

  sortRegionsByName(reg1, reg2) {
    const regionName1 = reg1.name.toLowerCase();
    const regionName2 = reg2.name.toLowerCase();

    if (regionName1 > regionName2) {
      return 1;
    } else if (regionName1 < regionName2) {
      return -1;
    }

    return 0;
  }

  createCluster() {
    if (!this.provisioning) {
      // compile frontend new-cluster model into api format
      const newClusterData: any = {};
      newClusterData.profile = {};

      newClusterData.cloudAccountName = this.selectedCloudAccount.name;
      newClusterData.clusterName = this.clusterName;
      newClusterData.profile.K8sVersion = this.clusterConfig.value.K8sVersion;
      newClusterData.profile.arch = this.clusterConfig.value.arch;
      newClusterData.profile.cidr = this.clusterConfig.value.cidr;
      newClusterData.profile.dockerVersion = this.clusterConfig.value.dockerVersion;
      newClusterData.profile.helmVersion = this.clusterConfig.value.helmVersion;
      newClusterData.profile.networkProvider = this.clusterConfig.value.networkProvider;
      newClusterData.profile.networkType = this.clusterConfig.value.networkType;
      newClusterData.profile.operatingSystem = this.clusterConfig.value.operatingSystem;
      newClusterData.profile.ubuntuVersion = this.clusterConfig.value.ubuntuVersion;
      newClusterData.profile.provider = this.selectedCloudAccount.provider;
      newClusterData.profile.rbacEnabled = true;
      newClusterData.profile.masterProfiles = this.nodesService.compileProfiles(this.selectedCloudAccount.provider, this.machines, 'Master');
      newClusterData.profile.nodesProfiles = this.nodesService.compileProfiles(this.selectedCloudAccount.provider, this.machines, 'Node');

      let region = this.availableRegions.filter(reg => {
        return reg.name === this.providerConfig.value.region
      })[0];
      newClusterData.profile.region = region.id;

      let exposedAddrObjects = new Array();
      for(let i=0; i<this.exposedAddressesArray.length; i++){
        let myObject = {"cidr" : this.exposedAddressesArray[i]};
        exposedAddrObjects.push(myObject);
      }
      newClusterData.profile.exposedAddresses = exposedAddrObjects;

      switch (newClusterData.profile.provider) {
        case 'aws':
          newClusterData.profile.cloudSpecificSettings = {
            aws_vpc_cidr: this.providerConfig.value.vpcCidr,
            aws_subnet_id: this.providerConfig.value.subnetId,
          };

          newClusterData.profile.publicKey = this.providerConfig.value.publicKey;
          break;
        case 'gce':
          newClusterData.profile.publicKey = this.providerConfig.value.publicKey;
          break;
        case 'azure':
          newClusterData.profile.cloudSpecificSettings = { azureVNetCIDR: this.providerConfig.value.azureVNetCIDR }
          newClusterData.profile.publicKey = this.providerConfig.value.publicKey;
          break;
      }

      this.provisioning = true;
      this.subscriptions.add(this.supergiant.Kubes.create(newClusterData).subscribe(
        data => {
          this.displaySuccess(newClusterData);
          this.router.navigate(['/clusters/', data.clusterId]);
          this.provisioning = false;
        },
        err => {
          this.displayError(err);
          this.provisioning = false;
        },
      ));
    }
  }

  displaySuccess(model) {
    this.notifications.display(
      'success',
      'Kube: ' + model.clusterName,
      'Created!',
    );
  }

  displayError(err) {
    let msg: string;

    if (err.error.userMessage) {
      msg = err.error.userMessage;
    } else {
      msg = err.error
    }

    this.notifications.display(
      'error',
      'Error: ',
      msg
    );
  }

  getAwsAvailabilityZones(region) {
    const accountName = this.selectedCloudAccount.name;

    return this.supergiant.CloudAccounts.getAwsAvailabilityZones(accountName, region.name);
  }

  getGCEAvailabilityZones(region) {
    const accountName = this.selectedCloudAccount.name;

    return this.supergiant.CloudAccounts.getGCEAvailabilityZones(accountName, region.name);
  }

  selectAz(zone, idx) {
    const accountName = this.selectedCloudAccount.name;
    const region = this.providerConfig.value.region;

    this.machinesLoading = true;

    switch (this.selectedCloudAccount.provider) {
      case 'aws':
        this.supergiant.CloudAccounts.getAwsMachineTypes(accountName, region, zone).subscribe(
          types => {
            this.machines[idx] = {
              ...this.machines[idx],
              availableMachineTypes: types.sort(),
            };

            this.machinesLoading = false;
          },
          err => {
            console.error(err);
            this.machinesLoading = false;
          },
        );
        break;
      case 'gce':
        this.supergiant.CloudAccounts.getGCEMachineTypes(accountName, region, zone).subscribe(
          types => {
            this.machines[idx] = {
              ...this.machines[idx],
              availableMachineTypes: types.sort(),
            };

            this.machinesLoading = false;
          },
          err => {
            console.error(err);
            this.machinesLoading = false;
          },
        );
        break;
    }
  }

  selectRegion(regionName) {

    let region = this.availableRegions.filter(reg => {
      return reg.name === regionName
    })[0];

    switch (this.selectedCloudAccount.provider) {
      case 'digitalocean':
        this.availableMachineTypes = sortDigitalOceanMachineTypes(region.AvailableSizes);

        if (this.machines.length === 0) {
          this.machines.push(BLANK_MACHINE_TEMPLATE);
        }
        break;
      case 'aws':
        this.azsLoading = true;
        this.getAwsAvailabilityZones(region).subscribe(
          azList => {
            this.availabilityZones = azList.sort();
            this.azsLoading = false;
          },
          err => {
            console.error(err);
            this.azsLoading = false;
          },
        );
        break;
      case 'gce':
        this.azsLoading = true;
        this.getGCEAvailabilityZones(region).subscribe(
          azList => {
            this.availabilityZones = azList.sort();
            this.azsLoading = false;
          },
          err => {
            console.error(err);
            this.azsLoading = false;
          },
        );
        break;
      case 'azure':
        // 'Basic_' VMs don't support load balancers
        const filterName = "Basic_";
        this.availableMachineTypes = region.AvailableSizes.filter(vmName => !vmName.includes(filterName));
    }
  }

  addBlankMachine() {
    const lastMachine = this.machines[this.machines.length - 1];

    this.machines.push(
      Object.assign({}, lastMachine),
    );
  }

  deleteMachine(idx) {
    if (this.machines.length === 1) {
      return;
    }

    this.machines.splice(idx, 1);
    this.validateMachineConfig();
  }

  selectCloudAccount(cloudAccount) {
    this.selectedCloudAccount = cloudAccount;

    this.availableRegions = null;
    this.availabilityZones = null;
    this.availableMachineTypes = null;
    this.availableRegionNames = new Array<string>();

    // TODO: quick fix to get pre-release cut
    // move to class and create new instance
    this.machines = [
      {
        machineType: null,
        qty: 1,
        availabilityZone: '',
        availableMachineTypes: null,
        role: 'Master',
        filter: '',
      },
      {
        machineType: null,
        qty: 1,
        availabilityZone: '',
        availableMachineTypes: null,
        role: 'Node',
        filter: '',
      },
    ];

    switch (this.selectedCloudAccount.provider) {
      case 'digitalocean':
        this.providerConfig = this.formBuilder.group({
          region: ['', Validators.required],
        });
        break;
      case 'aws':
        this.providerConfig = this.formBuilder.group({
          region: ['', Validators.required],
          vpcCidr: ['10.2.0.0/16', [Validators.required, this.validCidr()]],
          keypairName: [''],
          subnetId: [''],
          publicKey: ['', Validators.required],
        });
        break;
      case 'gce':
        this.providerConfig = this.formBuilder.group({
          region: ['', Validators.required],
          publicKey: ['', Validators.required],
        });
        break;
      case 'azure':
        this.providerConfig = this.formBuilder.group({
          region: ['', Validators.required],
          azureVNetCIDR: ['10.0.0.0/16', [Validators.required, this.validCidr()]],
          publicKey: ['', Validators.required]
        });
        break;
    }

    this.regionsLoading = true;
    this.subscriptions.add(this.supergiant.CloudAccounts.getRegions(cloudAccount.name).subscribe(
      regionList => {
        this.availableRegions = regionList.regions.sort(this.sortRegionsByName);
        this.availableRegionNames = this.availableRegions.map(n => n.name);
        this.regionsLoading = false;
      },
      err => {
        this.displayError(err);
        this.regionsLoading = false;
      },
    ));
  }

  validMachine(machine) {
    if (
      machine.machineType != null &&
      machine.role != null &&
      typeof (machine.qty) == 'number'
    ) {
      return true;
    } else {
      return false;
    }
    return false;
  }

  validateMachineConfig(currentMachine: IMachineType = null) {
    const selectedOption: MatOption = this.selectedMachineType.selected as MatOption;

    if (selectedOption) {
      const selectedMachineType: string = selectedOption.value;
    }

    if (this.machines.every(this.validMachine) && this.isOddNumberOfMasters() && this.hasMasterAndNode(this.machines)) {
      this.machinesConfigValid = true;
      this.displayMachinesConfigWarning = false;
    } else {
      this.machinesConfigValid = false;
    }
  }

  isOddNumberOfMasters() {
    const numberOfMasterProfiles = this.machines
      .filter(m => m.role === 'Master')
      .map(m => m.qty)
      .reduce((prev, next) => prev + next);
    return (numberOfMasterProfiles) % 2 !== 0;
  }

  hasMasterAndNode(machines) {
    const masters = machines.filter(m => m.role == MachineRoles.master)

    const nodes = machines.filter(m => m.role == MachineRoles.node)

    return (masters.length > 0 && nodes.length > 0);
  }

  nextStep() {

    if (this.stepper.selectedIndex === StepIndexes.Review) {
      this.validateMachineConfig();
    }

    if (this.machinesConfigValid) {
      this.stepper.next();
    } else {
      this.displayMachinesConfigWarning = true;
    }
  }

  uniqueClusterName(unavailableNames): ValidatorFn {
    return (name: AbstractControl): { [key: string]: any } | null => {
      return unavailableNames.has(name.value) ? { 'nonUniqueName': { value: name.value } } : null;
    };
  }

  validCidr(): ValidatorFn {
    return (userCidr: AbstractControl): { [key: string]: any } | null => {
      const validCidr = cidrRegex({ exact: true }).test(userCidr.value);
      return validCidr ? null : { 'invalidCidr': { value: userCidr.value } };
    };
  }

  allValidCidrs(): ValidatorFn {
    return (userCidrs: AbstractControl): { [key: string]: any } | null => {

      let cidrArray = this.toArray(userCidrs.value);
      let allValidCidrs = true;

      for(let i=0; i<cidrArray.length; i++){
        let validCidr = cidrRegex({ exact: true }).test(cidrArray[i]);
        if(!validCidr){
          allValidCidrs = false;
        }
      }

      return allValidCidrs ? null : { 'invalidCidrs': { value: true } };
    };
  }

  putExposedAddressesInArray(val) {
    this.exposedAddressesArray = this.toArray(val.target.value);
  }

  toArray(multiLineText : string) : Array<string> {
    return multiLineText.split("\n").filter(c => c != "");
  }

  showPublicKey() {
    this.initPublicKey(this.providerConfig.value.publicKey);
  }

  private initPublicKey(key) {
    const dialogRef = this.dialog.open(PublicKeyModalComponent, {
      width: '800px',
      data: { key: key }
    });

    return dialogRef;
  }

  allowSpaces(keyEvent){
    if(keyEvent.key==" "){
      keyEvent.stopPropagation();
    }
  }

}

import { Component, OnDestroy, OnInit, ViewEncapsulation, ViewChild } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ValidatorFn, AbstractControl } from '@angular/forms';
import { Router } from '@angular/router';
import { MatHorizontalStepper, MatOption, MatSelect } from '@angular/material';
import { Subscription, Observable } from 'rxjs';
import { Notifications } from '../../shared/notifications/notifications.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { NodeProfileService } from '../node-profile.service';
import { CLUSTER_OPTIONS } from './cluster-options.config';
import {
  DEFAULT_MACHINE_SET,
  BLANK_MACHINE_TEMPLATE,
} from 'app/clusters/new-cluster/new-cluster.component.config';
import { sortDigitalOceanMachineTypes } from 'app/clusters/new-cluster/new-cluster.helpers';
import { map } from 'rxjs/operators';
import { IMachineType } from './new-cluster.component.interface';

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
}

export enum StepIndexes {
  NameAndCloudAccount = 0,
  ClusterConfig = 1,
  ProvideConfig = 2,
  MachinesConfig = 3,
  Review = 4,
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
  nameAndCloudAccountForm: FormGroup;
  clusterConfig: FormGroup;
  providerConfig: FormGroup;
  unavailableClusterNames = new Set();
  machineTypesFilter = '';
  regionsFilter = '';

  @ViewChild(MatHorizontalStepper) stepper: MatHorizontalStepper;
  @ViewChild('selectedMachineType') selectedMachineType: MatSelect;

  constructor(
    private supergiant: Supergiant,
    private notifications: Notifications,
    private router: Router,
    private formBuilder: FormBuilder,
    private nodesService: NodeProfileService,
  ) {
  }

  ngOnInit() {
    this.getClusters();
    this.availableCloudAccounts$ = this.supergiant.CloudAccounts.get().pipe(
      map(cloudAccounts => cloudAccounts.sort()),
    );


    this.nameAndCloudAccountForm = this.formBuilder.group({
      name: ['', [
        Validators.required,
        this.uniqueClusterName(this.unavailableClusterNames),
        Validators.maxLength(12),
        Validators.pattern('^[A-Za-z]([-A-Za-z0-9\-]*[A-Za-z0-9\-])?$')]],
      cloudAccount: ['', Validators.required],
    });

    this.clusterConfig = this.formBuilder.group({
      K8sVersion: ['1.11.5', Validators.required],
      networkProvider: ['0.10.0', Validators.required],
      helmVersion: ['2.11.0', Validators.required],
      dockerVersion: ['17.06.0', Validators.required],
      ubuntuVersion: ['xenial', Validators.required],
      networkType: ['vxlan', Validators.required],
      cidr: ['10.0.0.0/16', [Validators.required, this.validCidr()]],
      operatingSystem: ['linux', Validators.required],
      arch: ['amd64', Validators.required],
    });

  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  get name() {
    return this.nameAndCloudAccountForm.get('name');
  }

  get cidr() {
    return this.clusterConfig.get('cidr');
  }

  get vpcCidr() {
    if (this.selectedCloudAccount && this.selectedCloudAccount.provider == 'aws') {
      return this.providerConfig.get('vpcCidr');
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
      newClusterData.profile = this.clusterConfig.value;

      newClusterData.cloudAccountName = this.selectedCloudAccount.name;
      newClusterData.clusterName = this.clusterName;
      newClusterData.profile.region = this.providerConfig.value.region.id;
      newClusterData.profile.provider = this.selectedCloudAccount.provider;
      newClusterData.profile.rbacEnabled = true;
      newClusterData.profile.masterProfiles = this.nodesService.compileProfiles(this.selectedCloudAccount.provider, this.machines, 'Master');
      newClusterData.profile.nodesProfiles = this.nodesService.compileProfiles(this.selectedCloudAccount.provider, this.machines, 'Node');

      switch (newClusterData.profile.provider) {
        case 'aws': {
          newClusterData.profile.cloudSpecificSettings = {
            aws_vpc_cidr: this.providerConfig.value.vpcCidr,
            aws_vpc_id: this.providerConfig.value.vpcId,
            aws_subnet_id: this.providerConfig.value.subnetId,
            aws_masters_secgroup_id: this.providerConfig.value.mastersSecurityGroupId,
            aws_nodes_secgroup_id: this.providerConfig.value.nodesSecurityGroupId,
          };

          newClusterData.profile.publicKey = this.providerConfig.value.publicKey;
        }
          break;
        case 'gce':

          newClusterData.profile.publicKey = this.providerConfig.value.publicKey;
          break;
      }

      this.provisioning = true;
      this.subscriptions.add(this.supergiant.Kubes.create(newClusterData).subscribe(
        (data: any) => {
          this.success(newClusterData);
          this.router.navigate(['/clusters/', data.clusterId]);
          this.provisioning = false;
        },
        (err) => {
          this.error(newClusterData, err);
          this.provisioning = false;
        },
      ));
    }
  }

  success(model) {
    this.notifications.display(
      'success',
      'Kube: ' + model.clusterName,
      'Created...',
    );
  }

  error(model, data) {
    this.notifications.display(
      'error',
      'Kube: ' + model.name,
      'Error:' + data.statusText);
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
    const region = this.providerConfig.value.region.name;

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

  selectRegion(region) {
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

    // TODO: quick fix to get pre-release cut
    // move to class and create new instance
    this.machines = [
      {
        machineType: null,
        qty: 1,
        availabilityZone: '',
        availableMachineTypes: null,
        role: 'Master',
      },
      {
        machineType: null,
        qty: 1,
        availabilityZone: '',
        availableMachineTypes: null,
        role: 'Node',
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
          vpcId: [''],
          vpcCidr: ['10.2.0.0/16', [Validators.required, this.validCidr()]],
          keypairName: [''],
          subnetId: [''],
          mastersSecurityGroupId: [''],
          nodesSecurityGroupId: [''],
          publicKey: ['', Validators.required],
        });
        break;
      case 'gce':
        this.providerConfig = this.formBuilder.group({
          region: ['', Validators.required],
          publicKey: ['', Validators.required],
        });
        break;
    }

    this.regionsLoading = true;
    this.subscriptions.add(this.supergiant.CloudAccounts.getRegions(cloudAccount.name).subscribe(
      regionList => {
        this.availableRegions = regionList;
        this.availableRegions.regions.sort(this.sortRegionsByName);
        this.regionsLoading = false;
      },
      err => {
        this.error({}, err);
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
      this.updateRecommendations(currentMachine, selectedMachineType);
    }

    if (this.machines.every(this.validMachine) && this.isOddNumberOfMasters()) {
      this.machinesConfigValid = true;
      this.displayMachinesConfigWarning = false;
    } else {
      this.machinesConfigValid = false;
    }
  }

  private updateRecommendations(currentMachine: IMachineType, selectedMachineType: string) {
    // checking recommendation for cluster size
    if (currentMachine && currentMachine.role === MachineRoles.master) {
      currentMachine.recommendedNodesCount =
        this.getRecommendedNodesCount(selectedMachineType) * currentMachine.qty;
    } else if (currentMachine) {
      currentMachine.recommendedNodesCount = 0;
    }
  }

  getRecommendedNodesCount(machineType: string): number {

    // TODO: move these consts into config file
    const AWS_RECOMENDATIONS = {
      'm3.medium': 5,
      'm3.large': 10,
      'm3.xlarge': 100,
      'm3.2xlarge': 250,
      'c4.4xlarge': 500,
      'c4.8xlarge': 500,
    };

    const DO_RECOMMENDATIONS = {
      's-1vcpu-3gb': 5,
      's-4vcpu-8gb': 10,
      's-6vcpu-16gb': 100,
      's-8vcpu-32gb': 250,
      's-16vcpu-64gb': 500,
    };

    const GCE_RECOMMENDATIONS = {
      'n1-standard-1': 5,
      'n1-standard-2': 10,
      'n1-standard-4': 100,
      'n1-standard-8': 250,
      'n1-standard-16': 500,
      'n1-standard-32': 1000,
    };


    switch (this.selectedCloudAccount.provider) {
      case CloudProviders.aws:
        return AWS_RECOMENDATIONS[machineType];
      case CloudProviders.digitalocean:
        return DO_RECOMMENDATIONS[machineType];
      case CloudProviders.gce:
        return GCE_RECOMMENDATIONS[machineType];

      default:
        return 0;
    }
  }

  isOddNumberOfMasters() {
    const numberOfMasterProfiles = this.machines
      .filter(m => m.role === 'Master')
      .map(m => m.qty)
      .reduce((prev, next) => prev + next);
    return (numberOfMasterProfiles) % 2 !== 0;
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

  machineTypesFilterCallback = (val) => {
    if (this.machineTypesFilter === '') {
      return val;
    }

    return val.toLowerCase().indexOf(this.machineTypesFilter.toLowerCase()) > -1;
  };

  regionsFilterCallback = (val) => {
    if (this.regionsFilter === '') {
      return val.name;
    }

    return val.name.toLowerCase().indexOf(this.regionsFilter.toLowerCase()) > -1;
  };
}

import { Component, OnDestroy, OnInit, ViewEncapsulation, ViewChild } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ValidatorFn, AbstractControl } from '@angular/forms';
import { Router }                 from '@angular/router';
import { MatHorizontalStepper }   from '@angular/material';
import { Subscription }           from 'rxjs';
import { Notifications }          from '../../shared/notifications/notifications.service';
import { Supergiant }             from '../../shared/supergiant/supergiant.service';
import { NodeProfileService }     from "../node-profile.service";
import { CLUSTER_OPTIONS }        from "./cluster-options.config";

// compiler hack
declare var require: any;
const cidrRegex = require('cidr-regex');

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
  availableRegions: any;
  selectedRegion: any;
  availableMachineTypes: Array<any>;
  regionsLoading = false;
  machinesLoading = false;

  // aws vars
  availabilityZones: Array<any>;
  azsLoading = false;

  machines = [{
    machineType: null,
    role: "Master",
    qty: 1
  }];

  clusterOptions = CLUSTER_OPTIONS;

  machinesConfigValid: boolean;
  displayMachinesConfigWarning: boolean;
  provisioning = false;
  nameAndCloudAccountConfig: FormGroup;
  clusterConfig: FormGroup;
  providerConfig: FormGroup;
  unavailableClusterNames = new Set();

  @ViewChild(MatHorizontalStepper) stepper: MatHorizontalStepper;

  constructor(
    private supergiant: Supergiant,
    private notifications: Notifications,
    private router: Router,
    private formBuilder: FormBuilder,
    private nodesService: NodeProfileService,
  ) { }


  getCloudAccounts() {
    this.subscriptions.add(this.supergiant.CloudAccounts.get().subscribe(
      (cloudAccounts) => {
        this.availableCloudAccounts = cloudAccounts.sort();
      })
    );
  }

  getClusters() {
    this.supergiant.Kubes.get().subscribe(
      clusters => clusters.map(c => this.unavailableClusterNames.add(c.name)),
      err => console.error(err)
    )
  }

  sortRegionsByName(reg1, reg2) {
    const regionName1 = reg1.name.toLowerCase();
    const regionName2 = reg2.name.toLowerCase();

    let comparison = 0;
    if (regionName1 > regionName2) {
      comparison = 1;
    } else if (regionName1 < regionName2) {
      comparison = -1;
    }
    return comparison;
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
      newClusterData.profile.masterProfiles = this.nodesService.compileProfiles(this.selectedCloudAccount.provider, this.machines, "Master");
      newClusterData.profile.nodesProfiles = this.nodesService.compileProfiles(this.selectedCloudAccount.provider, this.machines, "Node");

      switch (newClusterData.profile.provider) {
        case "aws": {
          newClusterData.profile.cloudSpecificSettings = {
            aws_az: this.providerConfig.value.availabilityZone,
            aws_vpc_cidr: this.providerConfig.value.vpcCidr,
            aws_vpc_id: this.providerConfig.value.vpcId,
            aws_keypair_name: this.providerConfig.value.keypairName,
            aws_subnet_id: this.providerConfig.value.subnetId,
            aws_masters_secgroup_id: this.providerConfig.value.mastersSecurityGroupId,
            aws_nodes_secgroup_id: this.providerConfig.value.nodesSecurityGroupId
          };

          newClusterData.profile.publicKey = this.providerConfig.value.publicKey;
        }
      }

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

  getAwsAvailabilityZones(region) {
    const accountName = this.selectedCloudAccount.name;

    return this.supergiant.CloudAccounts.getAwsAvailabilityZones(accountName, region.name);
  }

  selectAz(zone) {
    const accountName = this.selectedCloudAccount.name;
    const region = this.providerConfig.value.region.name;

    this.machinesLoading = true;
    this.supergiant.CloudAccounts.getAwsMachineTypes(accountName, region, zone).subscribe(
      types => {
        this.availableMachineTypes = types.sort();
        this.machinesLoading = false;
      },
      err => {
        console.error(err);
        this.machinesLoading = false;
      }
    )
  }

  selectRegion(region) {
    switch (this.selectedCloudAccount.provider) {
      case "digitalocean":
        this.availableMachineTypes = region.AvailableSizes.sort();
        if (this.machines.length === 0) {
          this.machines.push({
            machineType: null,
            role: null,
            qty: 1
          });
        }
        break;

      case "aws":
        this.azsLoading = true;
        this.getAwsAvailabilityZones(region).subscribe(
          azList => {
            this.availabilityZones = azList.sort();
            this.azsLoading = false
          },
          err => {
            console.error(err);
            this.azsLoading = false
          }
        );
        break;
    }
  }

  addBlankMachine(e?) {
    if (e) {
      if (e.keyCode === 13) {
        this.machines.push({
          machineType: null,
          role: null,
          qty: 1
        })
      }
    } else {
      this.machines.push({
        machineType: null,
        role: null,
        qty: 1
      })
    }
  }

  deleteMachine(idx, e?) {
    if(this.machines.length === 1) return;

    if (e) {
      if(e.keyCode === 13) {
        this.machines.splice(idx, 1);
        this.checkForValidMachinesConfig();
      }
    } else {
      this.machines.splice(idx, 1);
      this.checkForValidMachinesConfig();
    }
  }

  selectCloudAccount(cloudAccount) {
    this.selectedCloudAccount = cloudAccount;

    this.availableRegions = null;
    this.availabilityZones = null;
    this.availableMachineTypes = null;
    this.machines = [{
      machineType: null,
      role: "Master",
      qty: 1
    }];

    switch (this.selectedCloudAccount.provider) {
      case "digitalocean":
        this.providerConfig = this.formBuilder.group({
          region: ["", Validators.required]
        });
        break;

      case "aws":
        this.providerConfig = this.formBuilder.group({
          region: ["", Validators.required],
          availabilityZone: ["", Validators.required],
          vpcId: ["default", Validators.required],
          vpcCidr: ["10.2.0.0/16", [Validators.required, this.validCidr()]],
          keypairName: [""],
          subnetId: ["default", Validators.required],
          mastersSecurityGroupId: [""],
          nodesSecurityGroupId: [""],
          publicKey: ["", Validators.required]
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
      }
    ))
  }

  validMachine(machine) {
    if (
        machine.machineType != null &&
        machine.role != null &&
        typeof(machine.qty) == "number"
      ) {
      return true
   } else { return false }
  }

  checkForValidMachinesConfig() {
    if (this.machines.every(this.validMachine)) {
      this.machinesConfigValid = true;
      this.displayMachinesConfigWarning = false;
    } else {
      this.machinesConfigValid = false;
    }
  }

  machinesNext() {
    this.checkForValidMachinesConfig();

    if (this.machinesConfigValid) {
      this.stepper.next();
    } else {
      this.displayMachinesConfigWarning = true;
    }
  }

  uniqueClusterName(unavailableNames): ValidatorFn {
    return (name: AbstractControl): {[key: string]: any} | null => {
      return unavailableNames.has(name.value) ? {'nonUniqueName': {value: name.value}} : null;
    }
  }

  validCidr(): ValidatorFn {
    return (userCidr: AbstractControl): {[key: string]: any} | null => {
      const validCidr = cidrRegex({exact: true}).test(userCidr.value);
      return validCidr ? null : {"invalidCidr": {value: userCidr.value}};
    }
  }

  get name() { return this.nameAndCloudAccountConfig.get('name'); }

  get cidr() { return this.clusterConfig.get('cidr'); }

  get vpcCidr() {
    if (this.selectedCloudAccount && this.selectedCloudAccount.provider == "aws") {
      return this.providerConfig.get('vpcCidr');
    } else {return true}
  }

  ngOnInit() {
    this.getClusters();
    this.getCloudAccounts();

    this.nameAndCloudAccountConfig = this.formBuilder.group({
      name: ["", [
        Validators.required,
        this.uniqueClusterName(this.unavailableClusterNames),
        Validators.maxLength(12),
        Validators.pattern('([-A-Za-z0-9\-]*[A-Za-z0-9\-])?$')]],
      cloudAccount: ["", Validators.required]
    })

    this.clusterConfig = this.formBuilder.group({
      K8sVersion: ["1.11.1", Validators.required],
      flannelVersion: ["0.10.0", Validators.required],
      helmVersion: ["2.11.0", Validators.required],
      dockerVersion: ["17.06.0", Validators.required],
      ubuntuVersion: ["xenial", Validators.required],
      networkType: ["vxlan", Validators.required],
      cidr: ["10.0.0.0/16", [Validators.required, this.validCidr()]],
      operatingSystem: ["linux", Validators.required],
      arch: ["amd64", Validators.required],
      rbacEnabled: [false, Validators.required]
    });
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

}

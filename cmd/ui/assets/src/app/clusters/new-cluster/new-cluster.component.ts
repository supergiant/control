import { Component, OnDestroy, OnInit, ViewEncapsulation } from '@angular/core';
import { FormBuilder, FormGroup } from '@angular/forms';
import { Router }                 from '@angular/router';
import { Subscription }           from 'rxjs';
import { Notifications }          from '../../shared/notifications/notifications.service';
import { Supergiant }             from '../../shared/supergiant/supergiant.service';
import { NodeProfileService }     from "../node-profile.service";
import { CLUSTER_OPTIONS }        from "./cluster-options.config";

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
  availableMachineTypes: Array<any>;

  // aws vars
  availabilityZones: Array<any>;

  machines = [{
    machineType: null,
    role: "Master",
    qty: 1
  }];

  clusterOptions = CLUSTER_OPTIONS;

  provisioning = false;
  clusterConfig: FormGroup;
  providerConfig: FormGroup;

  constructor(
    private supergiant: Supergiant,
    private notifications: Notifications,
    private router: Router,
    private formBuilder: FormBuilder,
    private nodesService: NodeProfileService,
  ) {
  }


  getCloudAccounts() {
    this.subscriptions.add(this.supergiant.CloudAccounts.get().subscribe(
      (cloudAccounts) => {
        this.availableCloudAccounts = cloudAccounts;
      })
    );
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
          }
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

    this.supergiant.CloudAccounts.getAwsMachineTypes(accountName, region, zone).subscribe(
      types => this.availableMachineTypes = types.sort(),
      err => console.error(err)
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
        this.getAwsAvailabilityZones(region).subscribe(
          azList => this.availabilityZones = azList,
          err => console.error(err)
        );
        break;
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
    this.selectedCloudAccount = cloudAccount;

    switch (this.selectedCloudAccount.provider) {
      case "digitalocean":
        this.providerConfig = this.formBuilder.group({
          region: [""]
        });
        break;

      case "aws":
        this.providerConfig = this.formBuilder.group({
          region: [""],
          availabilityZone: [""],
          vpcId: ["default"],
          vpcCidr: ["10.2.0.0/16"],
          keypairName: [""],
          subnetId: ["default"],
          mastersSecurityGroupId: [""],
          nodesSecurityGroupId: [""]
        });
        break;
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
      helmVersion: ["2.11.0"],
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

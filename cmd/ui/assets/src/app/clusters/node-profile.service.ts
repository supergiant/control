import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class NodeProfileService {

  constructor() { }

  compileProfiles(provider, machines, role) {
    const filteredMachines = machines.filter(m => m.role == role);
    const compiledProfiles = [];

    switch (provider) {
      case "digitalocean":
        filteredMachines.forEach(m => {
          for (var i = 0; i < m.qty; i++) {
            compiledProfiles.push({ image: "ubuntu-16-04-x64", size: m.machineType })
          }
        });
        break;
      case "aws":
        filteredMachines.forEach(m => {
          for (var i = 0; i < m.qty; i++) {
            compiledProfiles.push({
              volumeSize: "80",
              size: m.machineType,
              ebsOptimized: "true",
              hasPublicAddr: "true"
            })
          }
        });
        break;
    }
    return compiledProfiles;
  }


  // selectAz(zone) {
  //   const accountName = this.selectedCloudAccount.name;
  //   const region = this.providerConfig.value.region.name;
  //
  //   this.supergiant.CloudAccounts.getAwsMachineTypes(accountName, region, zone).subscribe(
  //     types => this.availableMachineTypes = types.sort(),
  //     err => console.error(err)
  //   )
  // }

  // selectRegion(region) {
  //   switch (this.selectedCloudAccount.provider) {
  //     case "digitalocean":
  //       this.availableMachineTypes = region.AvailableSizes.sort();
  //       if (this.machines.length === 0) {
  //         this.machines.push({
  //           machineType: null,
  //           role: null,
  //           qty: 1
  //         });
  //       }
  //       break;
  //
  //     case "aws":
  //       this.getAwsAvailabilityZones(region).subscribe(
  //         azList => this.availabilityZones = azList,
  //         err => console.error(err)
  //       )
  //       break;
  //   }
  // }
}

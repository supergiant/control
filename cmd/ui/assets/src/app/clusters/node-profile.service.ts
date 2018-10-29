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

}

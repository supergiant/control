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
      case 'digitalocean':
        filteredMachines.forEach(m => {
          for (let i = 0; i < m.qty; i++) {
            compiledProfiles.push({ image: 'ubuntu-16-04-x64', size: m.machineType });
          }
        });
        break;
      case 'aws':
        filteredMachines.forEach(m => {
          for (let i = 0; i < m.qty; i++) {
            compiledProfiles.push({
              volumeSize: '80',
              size: m.machineType,
              ebsOptimized: 'true',
              hasPublicAddr: 'true',
              availabilityZone: m.availabilityZone
            });
          }
        });
        break;
      case 'gce':
        filteredMachines.forEach(m => {
          for (let i = 0; i < m.qty; i++) {
            compiledProfiles.push({
              size: m.machineType,
              availabilityZone: m.availabilityZone
            });
          }
        });
        break;
      case 'azure':
        filteredMachines.forEach(m => {
          for (let i = 0; i < m.qty; i++) {
            compiledProfiles.push({ 'vmSize': m.machineType });
          }
        });
        break;
    }
    return compiledProfiles;
  }
}

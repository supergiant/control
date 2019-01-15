import { getDigitalOceanMachineType } from 'app/clusters/new-cluster/new-cluster.helpers';

describe('NewClusterHelpersComponent', () => {

  describe('getDigitalOceanMachineType', () => {
    let machineSize: string;
    let expectedResutl: number;

    machineSize = '1mb';
    expectedResutl = 0;
    it(`should return ${expectedResutl} for ${machineSize}`, () => {
      const result = getDigitalOceanMachineType(machineSize);
      expect(result).toEqual(expectedResutl);
    });

    machineSize = '1gb';
    expectedResutl = 1;
    it(`should return ${expectedResutl} for ${machineSize}`, () => {
      const result = getDigitalOceanMachineType(machineSize);
      expect(result).toEqual(expectedResutl);
    });

    machineSize = 'm-4';
    expectedResutl = 2;
    it(`should return ${expectedResutl} for ${machineSize}`, () => {
      const result = getDigitalOceanMachineType(machineSize);
      expect(result).toEqual(expectedResutl);
    });

    machineSize = 'm-16gb';
    expectedResutl = 3;
    it(`should return ${expectedResutl} for ${machineSize}`, () => {
      const result = getDigitalOceanMachineType(machineSize);
      expect(result).toEqual(expectedResutl);
    });

    machineSize = 's-1cpu-1gb';
    expectedResutl = 4;
    it(`should return ${expectedResutl} for ${machineSize}`, () => {
      const result = getDigitalOceanMachineType(machineSize);
      expect(result).toEqual(expectedResutl);
    });

    machineSize = 's-12cpu-1gb';
    expectedResutl = 4;
    it(`should return ${expectedResutl} for ${machineSize}`, () => {
      const result = getDigitalOceanMachineType(machineSize);
      expect(result).toEqual(expectedResutl);
    });
  });
});

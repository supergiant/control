import {Injectable} from '@angular/core';
import {UtilService} from '../util/util.service';

@Injectable()
export class CloudAccounts {
  cloudAccountsPath = '/v1/api/accounts';

  constructor(private util: UtilService) { }

  public get(cloudAccountName?) {
    if (cloudAccountName) {
      return this.util.fetch(this.cloudAccountsPath + '/' + cloudAccountName);
    }
    return this.util.fetch(this.cloudAccountsPath);
  }

  public getRegions(cloudAccountName) {
    return this.util.fetch(this.cloudAccountsPath + '/' + cloudAccountName + '/' + 'regions');
  }

  public getAwsAvailabilityZones(cloudAccountName, region) {
    return this.util.fetch(this.cloudAccountsPath + '/' + cloudAccountName + '/' + 'regions' + '/' + region + '/az');
  }

  public getAwsMachineTypes(cloudAccountName, region, az) {
    return this.util.fetch(this.cloudAccountsPath + '/' + cloudAccountName + '/' + 'regions' + '/' + region + '/az/' + az + '/types');
  }
  public create(data) {
    return this.util.post(this.cloudAccountsPath, data);
  }

  public update(cloudAccountName, data) {
    return this.util.update(this.cloudAccountsPath + '/' + cloudAccountName, data);
  }

  public delete(cloudAccountName) {
    return this.util.destroy(this.cloudAccountsPath + '/' + cloudAccountName);
  }
}

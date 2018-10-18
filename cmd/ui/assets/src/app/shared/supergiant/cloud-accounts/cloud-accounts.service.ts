import { Injectable } from '@angular/core';
import { UtilService } from '../util/util.service';

@Injectable()
export class CloudAccounts {
  cloudAccountsPath = '/v1/api/accounts';

  constructor(private util: UtilService) { }
  public get(id?) {
    if (id) {
      return this.util.fetch(this.cloudAccountsPath + '/' + id);
    }
    return this.util.fetch(this.cloudAccountsPath);
  }
  public getRegions(id) {
    return this.util.fetch(this.cloudAccountsPath + '/' + id + '/' + 'regions');
  }
  public getAwsAvailabilityZones(id, region) {
    return this.util.fetch(this.cloudAccountsPath + '/' + id + '/' + 'regions' + '/' + region + '/az');
  }
  public getAwsMachineTypes(id, region, az) {
    return this.util.fetch(this.cloudAccountsPath + '/' + id + '/' + 'regions' + '/' + region + '/az/' + az + '/types');
  }
  public create(data) {
    return this.util.post(this.cloudAccountsPath, data);
  }
  public update(id, data) {
    return this.util.update(this.cloudAccountsPath + '/' + id, data);
  }
  public delete(id) {
    return this.util.destroy(this.cloudAccountsPath + '/' + id);
  }
}

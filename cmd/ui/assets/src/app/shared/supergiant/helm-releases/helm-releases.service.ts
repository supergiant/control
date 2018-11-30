import { Injectable } from '@angular/core';
import { UtilService } from '../util/util.service';

@Injectable()
export class HelmReleases {
  // TODO: everything about these paths need rethinking
  helmReleasesPath = '/v1/api';

  constructor(private util: UtilService) { }
  public get(clusterId, releaseName?) {
    if (releaseName) {
      return this.util.fetch(this.helmReleasesPath + '/kubes/' + clusterId + '/releases/' + releaseName);
    } else {
      return this.util.fetch(this.helmReleasesPath + '/kubes/' + clusterId + '/releases');
    }
  }
  public create(data) {
    return this.util.post(this.helmReleasesPath, data);
  }
  public update(clusterId, data) {
    return this.util.update(this.helmReleasesPath + '/' + clusterId, data);
  }
  public delete(releaseName, clusterId, purge) {
    return this.util.destroy(this.helmReleasesPath + '/kubes/' + clusterId + "/releases/" + releaseName + "?purge=" + purge);
  }
}

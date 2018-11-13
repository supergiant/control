import { Injectable } from '@angular/core';
import { UtilService } from '../util/util.service';

@Injectable()
export class HelmReleases {
  // TODO: everything about these paths need rethinking
  helmReleasesPath = '/v1/api';

  constructor(private util: UtilService) { }
  public get(id, releaseId?) {
    return this.util.fetch(this.helmReleasesPath + '/kubes/' + id + '/releases');
  }
  public create(data) {
    return this.util.post(this.helmReleasesPath, data);
  }
  public update(id, data) {
    return this.util.update(this.helmReleasesPath + '/' + id, data);
  }
  public delete(releaseName, kubeName, purge) {
    return this.util.destroy(this.helmReleasesPath + '/kubes/' + kubeName + "/releases/" + releaseName + "?purge=" + purge);
  }
}

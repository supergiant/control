import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/internal/Observable';
import { timeout } from 'rxjs/operators';
import { UtilService } from '../util/util.service';

@Injectable()
export class HelmReleases {
  // TODO: everything about these paths need rethinking
  helmReleasesPath = '/v1/api';

  constructor(private util: UtilService) {
  }

  public get(clusterId, releaseName?): Observable<any> {
    if (releaseName) {
      return this.util.fetch(this.helmReleasesPath + '/kubes/' + clusterId + '/releases/' + releaseName);
    }

    return this.util.fetch(this.helmReleasesPath + '/kubes/' + clusterId + '/releases')
      .pipe(
        timeout(1000),
      );
  }

  public create(data): Observable<any> {
    return this.util.post(this.helmReleasesPath, data);
  }

  public update(clusterId, data: any): Observable<any> {
    return this.util.update(this.helmReleasesPath + '/' + clusterId, data);
  }

  public delete(releaseName, clusterId, purge): Observable<any> {
    return this.util.destroy(this.helmReleasesPath + '/kubes/' + clusterId + '/releases/' + releaseName + '?purge=' + purge);
  }
}

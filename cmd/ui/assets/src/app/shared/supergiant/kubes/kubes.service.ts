import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { timeout } from 'rxjs/operators';
import { UtilService } from '../util/util.service';

@Injectable()
export class Kubes {
  kubesPath = '/v1/api/kubes';
  provisionPath = '/v1/api/provision';

  constructor(private util: UtilService) {
  }

  // TODO: add observable response models types
  public get(id?): Observable<any> {
    if (id) {
      return this.util.fetch(this.kubesPath + '/' + id);
    }
    return this.util.fetch(this.kubesPath);
  }

  public create(data): Observable<any> {
    return this.util.post(this.provisionPath, data);
  }

  public getClusterMetrics(id): Observable<any> {
    return this.util.fetch(this.kubesPath + '/' + id + '/metrics')
      .pipe(
        timeout(1000),
      );
  }

  public getMachineMetrics(id): Observable<any> {
    return this.util.fetch(this.kubesPath + '/' + id + '/nodes/metrics')
      .pipe(
        timeout(1000),
      );
  }

  public getClusterServices(id): Observable<any> {
    return this.util.fetch(this.kubesPath + '/' + id + '/services')
      .pipe(
        timeout(1000),
      );
  }

  // adding this back so I don't have to touch apps component right now
  public schema(data?): Observable<any> {
    return this.util.post(this.kubesPath, data);
  }

  public provision(id, data): Observable<any> {
    return this.util.post(this.kubesPath + '/' + id + '/provision', data);
  }

  public update(id, data): Observable<any> {
    return this.util.update(this.kubesPath + '/' + id, data);
  }

  public delete(id): Observable<any> {
    return this.util.destroy(this.kubesPath + '/' + id + '?force=true');
  }
}

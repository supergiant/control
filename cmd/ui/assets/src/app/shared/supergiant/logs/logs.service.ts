import { Injectable } from '@angular/core';
import { UtilService } from '../util/util.service';

@Injectable()
export class Logs {
  logsPath = '/api/v0/log';

  constructor(private util: UtilService) { }
  public get() {
    return this.util.fetchNoMap(this.logsPath);
  }
}

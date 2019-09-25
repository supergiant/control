import { Injectable } from '@angular/core';
import { UtilService } from '../util/util.service';

@Injectable()
export class Logs {
  logsPath = '/api/v1/log';

  constructor(private util: UtilService) { }
  public get() {
    return this.util.fetchResponse(this.logsPath);
  }
}

import { Injectable } from '@angular/core';
import { Http, Response, Headers } from '@angular/http';
import { CloudAccounts } from './cloud-accounts/cloud-accounts.service'
import { Sessions } from './sessions/sessions.service'
import { UtilService } from './util/util.service'

@Injectable()
export class Supergiant {
constructor(
  public CloudAccounts: CloudAccounts,
  public Sessions : Sessions,
) {}
}

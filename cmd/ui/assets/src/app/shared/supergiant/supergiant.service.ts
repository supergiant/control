import { Injectable } from '@angular/core';
import { Http, Response, Headers } from '@angular/http';
import { CloudAccount } from './cloud-accounts/cloud-accounts.service'
import { UtilService } from './util/util.service'

@Injectable()
export class Supergiant {
constructor(public CloudAccounts: CloudAccount) {}
}

@Injectable()
export class SupergiantService {
constructor(public CloudAccounts: CloudAccount) {}
}

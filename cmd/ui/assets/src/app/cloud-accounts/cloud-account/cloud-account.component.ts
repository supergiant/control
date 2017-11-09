import { Component, Input } from '@angular/core';
import { CloudAccountsService } from '../cloud-accounts.service';

@Component({
  selector: '[app-cloud-account]', // tslint:disable-line
  templateUrl: './cloud-account.component.html',
  styleUrls: ['./cloud-account.component.scss']
})
export class CloudAccountComponent {
  @Input() cloudAccount: any;
  private show: boolean;
  constructor(public cloudAccountsService: CloudAccountsService) { }


}

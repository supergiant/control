import { Component, Input } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { CloudAccountsService } from '../cloud-accounts.service';

@Component({
  selector: '[app-cloud-account]', // tslint:disable-line
  templateUrl: './cloud-account.component.html',
  styleUrls: ['./cloud-account.component.scss']
})
export class CloudAccountComponent {
  public subscriptions = new Subscription();
  @Input() cloudAccount: any;
  private show: boolean;
  public hasCloudAccount = false;

  constructor(
    public cloudAccountsService: CloudAccountsService,
    private supergiant: Supergiant,
  ) { }

  getCloudAccounts() {
    this.subscriptions.add(this.supergiant.CloudAccounts.get().subscribe(
      (cloudAccounts) => {
        if (Object.keys(cloudAccounts.items).length > 0) {
          this.hasCloudAccount = true;
        }
      })
    );
  }

  ngOnInit() {
    this.getCloudAccounts();
    console.log("Get Cloud");

    console.log(this.hasCloudAccount);
  }

}

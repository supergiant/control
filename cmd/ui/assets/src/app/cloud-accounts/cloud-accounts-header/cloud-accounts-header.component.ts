import { Component, OnInit } from '@angular/core';
import { CloudAccountsService } from '../cloud-accounts.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service'
import {CloudAccountsComponent} from '../cloud-accounts.component'

@Component({
  selector: 'app-cloud-accounts-header',
  templateUrl: './cloud-accounts-header.component.html',
  styleUrls: ['./cloud-accounts-header.component.css']
})
export class CloudAccountsHeaderComponent implements OnInit {

  constructor(
    private cloudAccountsService: CloudAccountsService,
    private cloudAccountsComponant: CloudAccountsComponent,
    private supergiant: Supergiant
    ) {}

  ngOnInit() {
  }



  sendOpen(message){
      this.cloudAccountsService.openNewCloudServiceModal(message);
  }
  deleteCloudAccount() {
    var selectedItems = this.cloudAccountsService.returnSelectedCloudAccount()

    for(let provider of selectedItems){
      this.supergiant.CloudAccounts.delete(provider.id).subscribe(
        (data) => {
          if (data.status >= 200 && data.status <= 299) {
            this.cloudAccountsService.showNotification("success", "Cloud Account: " + provider.name, "Deleted...")
            this.cloudAccountsComponant.getAccounts()
           }else{
            this.cloudAccountsService.showNotification("error", "Cloud Account: " + provider.name, "Error:" + data.statusText)}},
        (err) => {
          if (err) {
            this.cloudAccountsService.showNotification("error", "Cloud Account: " + provider.name, "Error:" + err)}},
      );
    }
  }

}

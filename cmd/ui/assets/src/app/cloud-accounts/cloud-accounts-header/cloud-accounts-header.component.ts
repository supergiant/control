import { Component } from '@angular/core';
import { CloudAccountsService } from '../cloud-accounts.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service'
import {CloudAccountsComponent} from '../cloud-accounts.component'
import { Subscription } from 'rxjs/Subscription';
import { Notifications } from '../../shared/notifications/notifications.service'
import { SystemModalService } from '../../shared/system-modal/system-modal.service'

@Component({
  selector: 'app-cloud-accounts-header',
  templateUrl: './cloud-accounts-header.component.html',
  styleUrls: ['./cloud-accounts-header.component.css']
})
export class CloudAccountsHeaderComponent {
  private cloudAccountsSub: Subscription;
  providersObj: any;

  constructor(
    private cloudAccountsService: CloudAccountsService,
    private cloudAccountsComponant: CloudAccountsComponent,
    private supergiant: Supergiant,
    private notifications: Notifications,
    private systemModalService: SystemModalService,
    ) {}

  // After init, grab the schema
  ngAfterViewInit() {
    this.cloudAccountsSub = this.supergiant.CloudAccounts.schema().subscribe(
      (data) => { this.providersObj = data.json()},
      (err) => {this.notifications.display("warn", "Connection Issue.", err)});
  }

  // If new button if hit, the New dropdown is triggered.
  sendOpen(message){
      this.cloudAccountsService.openNewCloudServiceDropdownModal(message);
  }

  openSystemModal(message){
      this.systemModalService.openSystemModal(message);
  }
  // If the edit button is hit, the Edit modal is opened.
  editCloudAccount() {
    var selectedItems = this.cloudAccountsService.returnSelectedCloudAccount()

    if (selectedItems.length === 0) {
      this.notifications.display("warn", "Warning:", "No Provider Selected.")
    } else if (selectedItems.length > 1) {
      this.notifications.display("warn", "Warning:", "You cannot edit more than one provider at a time.")
    } else {
      this.providersObj.providers[selectedItems[0].provider].model = selectedItems[0]
      this.cloudAccountsService.openNewCloudServiceEditModal("Edit", selectedItems[0].provider, this.providersObj);
    }
  }

  // If the delete button is hit, the seleted accounts are deleted.
  deleteCloudAccount() {
    var selectedItems = this.cloudAccountsService.returnSelectedCloudAccount()
    if (selectedItems.length === 0) {
      this.notifications.display("warn", "Warning:", "No Provider Selected.")
    } else if (selectedItems.length > 1) {
      this.notifications.display("warn", "Warning:", "You cannot edit more than one provider at a time.")
    } else {
    for(let provider of selectedItems){
      this.supergiant.CloudAccounts.delete(provider.id).subscribe(
        (data) => {
          if (data.status >= 200 && data.status <= 299) {
            this.notifications.display("success", "Cloud Account: " + provider.name, "Deleted...")
            this.cloudAccountsComponant.getAccounts()
           }else{
            this.notifications.display("error", "Cloud Account: " + provider.name, "Error:" + data.statusText)}},
        (err) => {
          if (err) {
            this.notifications.display("error", "Cloud Account: " + provider.name, "Error:" + err)}},
      );
    }
  }
  }
}

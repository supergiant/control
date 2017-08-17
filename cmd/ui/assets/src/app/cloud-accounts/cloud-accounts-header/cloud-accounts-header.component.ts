import { Component, OnDestroy } from '@angular/core';
import { CloudAccountsService } from '../cloud-accounts.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service'
import {CloudAccountsComponent} from '../cloud-accounts.component'
import { Subscription } from 'rxjs/Subscription';
import { Notifications } from '../../shared/notifications/notifications.service'
import { SystemModalService } from '../../shared/system-modal/system-modal.service'
import { DropdownModalService } from '../../shared/dropdown-modal/dropdown-modal.service'
import { EditModalService } from '../../shared/edit-modal/edit-modal.service'
import { LoginComponent } from '../../login/login.component';
import { CloudAccountModel } from '../cloud-accounts.model';


@Component({
  selector: 'app-cloud-accounts-header',
  templateUrl: './cloud-accounts-header.component.html',
  styleUrls: ['./cloud-accounts-header.component.css']
})
export class CloudAccountsHeaderComponent {
  providersObj = new CloudAccountModel
  subscriptions = [];
  editID: number;

  constructor(
    private cloudAccountsService: CloudAccountsService,
    private cloudAccountsComponent: CloudAccountsComponent,
    private supergiant: Supergiant,
    private notifications: Notifications,
    private systemModalService: SystemModalService,
    private dropdownModalService: DropdownModalService,
    private editModalService: EditModalService,
    private loginComponent: LoginComponent,
    ) {}

    ngOnDestroy(){
      for (let subscription of this.subscriptions)  {
        subscription.unsubscribe();
      }
    }

  ngAfterViewInit() {
    this.subscriptions["dropdown"] = this.dropdownModalService.dropdownModalResponse.subscribe(
      (option) => {
        this.editModalService.open("Save", option, this.providersObj)},
      );

    this.subscriptions["edit"] = this.editModalService.editModalResponse.subscribe(
      (userInput) => {
        var action = userInput[0]
        var providerID = 1
        var model = userInput[2]

        if (action === "Edit") {
          this.supergiant.CloudAccounts.update(providerID, model).subscribe(
            (data) => {
              this.success(model)
              this.cloudAccountsComponent.getAccounts()},
            (err) => {this.error(model, err)},);
        } else {
          this.supergiant.CloudAccounts.create(model).subscribe(
            (data) => {
              this.success(model)
              this.cloudAccountsComponent.getAccounts()},
            (err) => {this.error(model, err)});
        }
      });
  }

  success(model){
    this.notifications.display(
      "success",
      "Cloud Account: " + model.name,
      "Created...",
    )
  }

  error(model, data) {
    this.notifications.display(
      "error",
      "Cloud Account: " + model.name,
      "Error:" + data.statusText)
  }
  // If new button if hit, the New dropdown is triggered.
  sendOpen(message){
     let providers = [];

     // Push available providers to an array. Displayed in the dropdown.
     for(let key in this.providersObj.providers){
       providers.push(key)
     }

      // Open Dropdown Modal
      this.dropdownModalService.open("New Cloud Account", "Cloud Account", providers)

  }

  openSystemModal(message){
      this.systemModalService.openSystemModal(message);
  }
  // If the edit button is hit, the Edit modal is opened.
  editCloudAccount() {
    var selectedItems = this.cloudAccountsService.returnSelected()
    var itemindex: string
    if (selectedItems.length === 0) {
      this.notifications.display("warn", "Warning:", "No Cloud Account Selected.")
    } else if (selectedItems.length > 1) {
      this.notifications.display("warn", "Warning:", "You cannot edit more than one Cloud Account at a time.")
    } else {
      for (let aprovider in this.providersObj.providers) {
        if (this.providersObj.providers[aprovider]["model"]["provider"] == selectedItems[0].provider) {
          this.providersObj.providers[aprovider]["model"] = selectedItems[0]
          itemindex = aprovider
        }
      }
      this.editModalService.open("Edit", itemindex, this.providersObj.providers);
    }
  }

  // If the delete button is hit, the seleted accounts are deleted.
  deleteCloudAccount() {
    var selectedItems = this.cloudAccountsService.returnSelected()
    if (selectedItems.length === 0) {
      this.notifications.display("warn", "Warning:", "No Cloud Account Selected.")
    } else {
    for(let provider of selectedItems){
      this.supergiant.CloudAccounts.delete(provider.id).subscribe(
        (data) => {
          if (data.status >= 200 && data.status <= 299) {
            this.notifications.display("success", "Cloud Account: " + provider.name, "Deleted...")
            this.cloudAccountsComponent.getAccounts()
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

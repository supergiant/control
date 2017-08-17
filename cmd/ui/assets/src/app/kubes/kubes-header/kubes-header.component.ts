import { Component } from '@angular/core';
import { KubesService } from '../kubes.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service'
import {KubesComponent} from '../kubes.component'
import { Subscription } from 'rxjs/Subscription';
import { Notifications } from '../../shared/notifications/notifications.service'
import { SystemModalService } from '../../shared/system-modal/system-modal.service'
import { DropdownModalService } from '../../shared/dropdown-modal/dropdown-modal.service'
import { EditModalService } from '../../shared/edit-modal/edit-modal.service'
import { LoginComponent } from '../../login/login.component';
import { KubesModel } from '../kubes.model';

@Component({
  selector: 'app-kubes-header',
  templateUrl: './kubes-header.component.html',
  styleUrls: ['./kubes-header.component.css']
})
export class KubesHeaderComponent {
  providersObj = {"providers": []}
  public subscription: any = {}

  kubesModel = new KubesModel
  constructor(
    private kubesService: KubesService,
    private kubesComponent: KubesComponent,
    private supergiant: Supergiant,
    private notifications: Notifications,
    private systemModalService: SystemModalService,
    private dropdownModalService: DropdownModalService,
    private editModalService: EditModalService,
    private loginComponent: LoginComponent
    ) {}

  // After init, grab the schema
  ngAfterViewInit() {
    let providers = [];
    let cloudAccountsList
    //this.providersObj = KubesModel
    this.subscription["accounts"]=this.supergiant.CloudAccounts.get().subscribe(
      (data) => { cloudAccountsList = data
      console.log(cloudAccountsList["items"])
      for(let account of cloudAccountsList["items"]){
           console.log(account)
           providers.push(account["name"])
           console.log(this.providersObj.providers)
           this.providersObj.providers[account["name"]] = this.kubesModel.providers[account["provider"]]
           console.log(this.providersObj.providers)
           console.log(this.providersObj.providers[account["name"]])
           this.providersObj.providers[account["name"]]["cloud_account_name"] = account["name"]
         }}
    )
  }

  // If new button if hit, the New dropdown is triggered.
  sendOpen(message){
     let providers = [];
     // Push available providers to an array. Displayed in the dropdown.
     for(let key in this.providersObj.providers){
       providers.push(key)
     }

      // Open Dropdown Modal
      this.subscription["button"] = this.dropdownModalService.open(
        "New Kube", "Providers", providers).subscribe(
          (option) => {
            this.subscription["edit_button"] = this.editModalService.open("Save", option, this.providersObj).subscribe(
              (userInput) => {
                var action = userInput[0]
                var providerID = 1
                var model = userInput[2]
                if (action === "Edit") {
                this.subscription["update_button"] = this.supergiant.Kubes.update(providerID, model).subscribe(
                  (data) => {
                    if (data.status >= 200 && data.status <= 299) {
                      this.notifications.display(
                        "success",
                        "Kube: " + model.name,
                        "Created...",
                      )
                      this.kubesComponent.getAccounts()
                    }else{
                      this.notifications.display(
                        "error",
                        "Kube: " + model.name,
                        "Error:" + data.statusText)
                      }},
                  (err) => {
                    if (err) {
                      this.notifications.display(
                        "error",
                        "Kube: " + model.name,
                        "Error:" + err)
                      }});
              } else {
                this.subscription["create_button"] = this.supergiant.Kubes.create(model).subscribe(
                  (data) => {
                    if (data.status >= 200 && data.status <= 299) {
                      this.notifications.display(
                        "success",
                        "Kube: " + model.name.name,
                        "Created...",
                      )
                      this.kubesComponent.getAccounts()
                    }else{
                      this.notifications.display(
                        "error",
                        "Kube: " + model.name.name,
                        "Error:" + data.statusText)
                      }},
                  (err) => {
                    if (err) {
                      this.notifications.display(
                        "error",
                        "Kube: " + model.name.name,
                        "Error:" + err)
                      }});}
              });
          });
      console.log(this.subscription)
      for (let sub of this.subscription){
        sub.unsubscribe()
        console.log("unsub")
      }
  }

  openSystemModal(message){
      this.systemModalService.openSystemModal(message);
  }
  // If the edit button is hit, the Edit modal is opened.
  editKube() {
    var selectedItems = this.kubesService.returnSelected()

    if (selectedItems.length === 0) {
      this.notifications.display("warn", "Warning:", "No Provider Selected.")
    } else if (selectedItems.length > 1) {
      this.notifications.display("warn", "Warning:", "You cannot edit more than one provider at a time.")
    } else {
      this.providersObj.providers[selectedItems[0].provider].model = selectedItems[0]
      this.editModalService.open("Edit", selectedItems[0].provider, this.providersObj);
    }
  }

  // If the delete button is hit, the seleted accounts are deleted.
  deleteCloudAccount() {
    var selectedItems = this.kubesService.returnSelected()
    if (selectedItems.length === 0) {
      this.notifications.display("warn", "Warning:", "No Provider Selected.")
    } else if (selectedItems.length > 1) {
      this.notifications.display("warn", "Warning:", "You cannot edit more than one provider at a time.")
    } else {
    for(let provider of selectedItems){
      this.supergiant.CloudAccounts.delete(provider.id).subscribe(
        (data) => {
          if (data.status >= 200 && data.status <= 299) {
            this.notifications.display("success", "Kube: " + provider.name, "Deleted...")
            this.kubesComponent.getAccounts()
           }else{
            this.notifications.display("error", "Kube: " + provider.name, "Error:" + data.statusText)}},
        (err) => {
          if (err) {
            this.notifications.display("error", "Kube: " + provider.name, "Error:" + err)}},
      );
    }
  }
  }
}

import { Component } from '@angular/core';
import { KubesService } from '../kubes.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service'
import {KubesComponent} from '../kubes.component'
import { Subscription } from 'rxjs/Subscription';
import { Notifications } from '../../shared/notifications/notifications.service'
import { SystemModalService } from '../../shared/system-modal/system-modal.service'
import { DropdownModalService } from '../../shared/dropdown-modal/dropdown-modal.service'
import { EditModalService } from '../../shared/edit-modal/edit-modal.service'


@Component({
  selector: 'app-kubes-header',
  templateUrl: './kubes-header.component.html',
  styleUrls: ['./kubes-header.component.css']
})
export class KubesHeaderComponent {
  providersObj: any;

  constructor(
    private kubesService: KubesService,
    private kubesComponant: KubesComponent,
    private supergiant: Supergiant,
    private notifications: Notifications,
    private systemModalService: SystemModalService,
    private dropdownModalService: DropdownModalService,
    private editModalService: EditModalService,
    ) {}

  // After init, grab the schema
  ngAfterViewInit() {
    this.supergiant.Kubes.schema().subscribe(
      (data) => { this.providersObj = data.json()},
      (err) => {this.notifications.display("warn", "Connection Issue.", err)});
  }

  // If new button if hit, the New dropdown is triggered.
  sendOpen(message){
     let providers = [];
     // Fetch options.
     this.supergiant.Kubes.schema().subscribe(
       (data) => { this.providersObj = data.json()
         // Push available providers to an array. Displayed in the dropdown.
         for(let key in this.providersObj.providers){
           providers.push(key)
         }
       });

      // Open Dropdown Modal
      this.dropdownModalService.open(
        "New Kube", "Providers", providers).subscribe(
          (option) => {
            this.editModalService.open("Save", option, this.providersObj).subscribe(
              (userInput) => {
                var action = userInput[0]
                var providerID = 1
                var model = userInput[2]
                if (action === "Edit") {
                this.supergiant.Kubes.update(providerID, model).subscribe(
                  (data) => {
                    if (data.status >= 200 && data.status <= 299) {
                      this.notifications.display(
                        "success",
                        "Kube: " + model.name,
                        "Created...",
                      )
                      this.kubesComponant.getAccounts()
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
                this.supergiant.Kubes.create(model).subscribe(
                  (data) => {
                    if (data.status >= 200 && data.status <= 299) {
                      this.notifications.display(
                        "success",
                        "Kube: " + model.name.name,
                        "Created...",
                      )
                      this.kubesComponant.getAccounts()
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
            this.kubesComponant.getAccounts()
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

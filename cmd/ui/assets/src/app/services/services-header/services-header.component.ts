import { Component } from '@angular/core';
import { ServicesService } from '../services.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service'
import { ServicesComponent } from '../services.component'
import { Subscription } from 'rxjs/Subscription';
import { Notifications } from '../../shared/notifications/notifications.service'
import { SystemModalService } from '../../shared/system-modal/system-modal.service'
import { DropdownModalService } from '../../shared/dropdown-modal/dropdown-modal.service'
import { EditModalService } from '../../shared/edit-modal/edit-modal.service'
import { LoginComponent } from '../../login/login.component';
import { ServicesModel } from '../services.model';

@Component({
  selector: 'app-services-header',
  templateUrl: './services-header.component.html',
  styleUrls: ['./services-header.component.css']
})
export class ServicesHeaderComponent {
  providersObj: any;
  subscriptions = new Subscription();

  constructor(
    private servicesService: ServicesService,
    private servicesComponent: ServicesComponent,
    private supergiant: Supergiant,
    private notifications: Notifications,
    private systemModalService: SystemModalService,
    private dropdownModalService: DropdownModalService,
    private editModalService: EditModalService,
    private loginComponent: LoginComponent,
  ) { }

  ngOnDestroy() {
    this.subscriptions.unsubscribe()
  }
  // After init, grab the schema
  ngAfterViewInit() {
    this.providersObj = ServicesModel
    this.subscriptions.add(this.dropdownModalService.dropdownModalResponse.subscribe(
      (option) => { this.editModalService.open("Save", option, this.providersObj) }, ))


    this.subscriptions.add(this.editModalService.editModalResponse.subscribe(
      (userInput) => {
        var action = userInput[0]
        var providerID = 1
        var model = userInput[2]
        if (action === "Edit") {
          this.subscriptions.add(this.supergiant.Nodes.update(providerID, model).subscribe(
            (data) => {
              this.success(model)
              this.servicesComponent.getAccounts()
            },
            (err) => { this.error(model, err) }))
        } else {
          this.subscriptions.add(this.supergiant.KubeResources.create(model).subscribe(
            (data) => {
              this.success(model)
              this.servicesComponent.getAccounts()
            },
            (err) => { this.error(model, err) }))
        }
      }))
  }

  success(model) {
    this.notifications.display(
      "success",
      "Service: " + model.name,
      "Created...",
    )
  }

  error(model, data) {
    this.notifications.display(
      "error",
      "Service: " + model.name,
      "Error:" + data.statusText)
  }

  // If new button if hit, the New dropdown is triggered.
  sendOpen(message) {
    let providers = [];
    // Push available providers to an array. Displayed in the dropdown.
    for (let key in this.providersObj.providers) {
      providers.push(key)
    }

    // Open Dropdown Modal
    this.dropdownModalService.open(
      "New Service", "Providers", providers)

  }

  openSystemModal(message) {
    this.systemModalService.openSystemModal(message);
  }
  // If the edit button is hit, the Edit modal is opened.
  editUser() {
    var selectedItems = this.servicesService.returnSelected()

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
    var selectedItems = this.servicesService.returnSelected()
    if (selectedItems.length === 0) {
      this.notifications.display("warn", "Warning:", "No Service Selected.")
    } else {
      for (let provider of selectedItems) {
        this.subscriptions.add(this.supergiant.CloudAccounts.delete(provider.id).subscribe(
          (data) => {
            this.notifications.display("success", "Service: " + provider.name, "Deleted...")
            this.servicesComponent.getAccounts()
          },
          (err) => {
            if (err) {
              this.notifications.display("error", "Service: " + provider.name, "Error:" + err)
            }
          },
        ))
      }
    }
  }
}

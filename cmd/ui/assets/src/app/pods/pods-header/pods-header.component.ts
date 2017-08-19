import { Component } from '@angular/core';
import { PodsService } from '../pods.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service'
import { PodsComponent } from '../pods.component'
import { Subscription } from 'rxjs/Subscription';
import { Notifications } from '../../shared/notifications/notifications.service'
import { SystemModalService } from '../../shared/system-modal/system-modal.service'
import { DropdownModalService } from '../../shared/dropdown-modal/dropdown-modal.service'
import { EditModalService } from '../../shared/edit-modal/edit-modal.service'
import { LoginComponent } from '../../login/login.component';
import { PodsModel } from '../pods.model';

@Component({
  selector: 'app-pods-header',
  templateUrl: './pods-header.component.html',
  styleUrls: ['./pods-header.component.css']
})
export class PodsHeaderComponent {
  providersObj: any;
  subscriptions = new Subscription();

  constructor(
    private podsService: PodsService,
    private podsComponent: PodsComponent,
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
    this.providersObj = PodsModel
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
              this.podsComponent.getAccounts()
            },
            (err) => { this.error(model, err) }))
        } else {
          this.subscriptions.add(this.supergiant.Nodes.create(model).subscribe(
            (data) => {
              this.success(model)
              this.podsComponent.getAccounts()
            },
            (err) => { this.error(model, err) }))
        }
      }))
  }

  success(model) {
    this.notifications.display(
      "success",
      "Pod: " + model.name,
      "Created...",
    )
  }

  error(model, data) {
    this.notifications.display(
      "error",
      "Pod: " + model.name,
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
      "New Pod", "Providers", providers)

  }

  openSystemModal(message) {
    this.systemModalService.openSystemModal(message);
  }
  // If the edit button is hit, the Edit modal is opened.
  editUser() {
    var selectedItems = this.podsService.returnSelected()

    if (selectedItems.length === 0) {
      this.notifications.display("warn", "Warning:", "No Pod Selected.")
    } else if (selectedItems.length > 1) {
      this.notifications.display("warn", "Warning:", "You cannot edit more than one Pod at a time.")
    } else {
      this.providersObj.providers[selectedItems[0].provider].model = selectedItems[0]
      this.editModalService.open("Edit", selectedItems[0].provider, this.providersObj);
    }
  }

  // If the delete button is hit, the seleted accounts are deleted.
  deleteCloudAccount() {
    var selectedItems = this.podsService.returnSelected()
    if (selectedItems.length === 0) {
      this.notifications.display("warn", "Warning:", "No Pod Selected.")
    } else {
      for (let provider of selectedItems) {
        this.subscriptions.add(this.supergiant.CloudAccounts.delete(provider.id).subscribe(
          (data) => {
            this.notifications.display("success", "Pod: " + provider.name, "Deleted...")
            this.podsComponent.getAccounts()
          },
          (err) => {
            this.notifications.display("error", "Pod: " + provider.name, "Error:" + err)
          },
        ))
      }
    }
  }
}

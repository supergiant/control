import { Component } from '@angular/core';
import { KubesService } from '../kubes.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service'
import { KubesComponent } from '../kubes.component'
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
  subscriptions = new Subscription();
  cloudAccountsList = [];

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
  ) { }

  ngOnDestroy() {
    this.subscriptions.unsubscribe()
  }
  // After init, grab the schema
  ngAfterViewInit() {
    //this.providersObj = KubesModel
    this.subscriptions.add(this.supergiant.CloudAccounts.get().subscribe(
      (data) => { this.cloudAccountsList = data.items }
    ))

    this.subscriptions.add(this.dropdownModalService.dropdownModalResponse.subscribe(
      (option) => {
        if (option != "closed") {
          let cloudAccount = this.cloudAccountsList.filter(resource => resource.name == option)[0]
          this.kubesModel.providers[cloudAccount.provider].model.cloud_account_name = cloudAccount.name
          this.editModalService.open("Save", cloudAccount.provider, this.kubesModel.providers)
        }
      },
      (err) => { console.log("ERROR: " + err) }
    ))

    this.subscriptions.add(this.editModalService.editModalResponse.subscribe(
      (userInput) => {
        if (userInput != "closed") {
          var action = userInput[0]
          var providerID = 1
          var model = userInput[2]
          if (action === "Edit") {
            this.subscriptions.add(this.supergiant.Kubes.update(providerID, model).subscribe(
              (data) => {
                console.log(data)
                this.success(model)
                this.kubesComponent.getAccounts()
              },
              (err) => { this.error(model, err) }))
          } else {
            this.subscriptions.add(this.supergiant.Kubes.create(model).subscribe(
              (data) => {
                this.success(model)
                this.kubesComponent.getAccounts()
              },
              (err) => { this.error(model, err) }))
          }
        }
      }))
  }

  success(model) {
    this.notifications.display(
      "success",
      "Kube: " + model.name,
      "Created...",
    )
  }

  error(model, data) {
    this.notifications.display(
      "error",
      "Kube: " + model.name,
      "Error:" + data.statusText)
  }

  sendOpen(message) {
    let providers = [];
    providers = this.cloudAccountsList.map((resource) => { return resource.name })
    this.dropdownModalService.open("New Kube", "Providers", providers)
  }

  openSystemModal(message) {
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
      this.kubesModel.providers[selectedItems[0].provider].model = selectedItems[0]
      this.editModalService.open("Edit", selectedItems[0].provider, this.kubesModel);
    }
  }

  // If the delete button is hit, the seleted accounts are deleted.
  deleteKube() {
    var selectedItems = this.kubesService.returnSelected()
    if (selectedItems.length === 0) {
      this.notifications.display("warn", "Warning:", "No Kube Selected.")
    } else {
      for (let provider of selectedItems) {
        this.subscriptions.add(this.supergiant.Kubes.delete(provider.id).subscribe(
          (data) => {
            this.notifications.display("success", "Kube: " + provider.name, "Deleted...")
            this.kubesService.resetSelected()
            this.kubesComponent.getAccounts()
          },
          (err) => {
            this.notifications.display("error", "Kube: " + provider.name, "Error:" + err)
          },
        ))
      }
    }
  }
}

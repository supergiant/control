import { Component, OnDestroy, AfterViewInit } from '@angular/core';
import { CloudAccountsService } from '../cloud-accounts.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Subscription } from 'rxjs/Subscription';
import { Notifications } from '../../shared/notifications/notifications.service';
import { SystemModalService } from '../../shared/system-modal/system-modal.service';
import { DropdownModalService } from '../../shared/dropdown-modal/dropdown-modal.service';
import { EditModalService } from '../../shared/edit-modal/edit-modal.service';
import { LoginComponent } from '../../login/login.component';
import { CloudAccountModel } from '../cloud-accounts.model';


@Component({
  selector: 'app-cloud-accounts-header',
  templateUrl: './cloud-accounts-header.component.html',
  styleUrls: ['./cloud-accounts-header.component.css']
})
export class CloudAccountsHeaderComponent implements OnDestroy, AfterViewInit {
  providersObj = new CloudAccountModel;
  subscriptions = new Subscription();
  editID: number;

  constructor(
    private cloudAccountsService: CloudAccountsService,
    private supergiant: Supergiant,
    private notifications: Notifications,
    private systemModalService: SystemModalService,
    private dropdownModalService: DropdownModalService,
    private editModalService: EditModalService,
    private loginComponent: LoginComponent,
  ) { }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  ngAfterViewInit() {
    this.subscriptions.add(this.dropdownModalService.dropdownModalResponse.subscribe(
      (option) => {
        if (option !== 'closed') {
          this.editModalService.open('Save', option, this.providersObj.providers);
        }
      },
    ));

    this.subscriptions.add(this.editModalService.editModalResponse.subscribe(
      (userInput) => {
        if (userInput !== 'closed') {
          const action = userInput[0];
          const providerID = userInput[1];
          const model = userInput[2];

          if (action === 'Edit') {
            this.supergiant.CloudAccounts.update(providerID, model).subscribe(
              (data) => {
                this.success(model);
              },
              (err) => { this.error(model, err); }, );
          } else {
            this.supergiant.CloudAccounts.create(model).subscribe(
              (data) => {
                this.success(model);
              },
              (err) => { this.error(model, err); });
          }
        }
      }
    ));
  }

  success(model) {
    this.notifications.display(
      'success',
      'Cloud Account: ' + model.name,
      'Created...',
    );
  }

  error(model, data) {
    this.notifications.display(
      'error',
      'Cloud Account: ' + model.name,
      'Error:' + data.statusText);
  }

  // If new button if hit, the New dropdown is triggered.
  sendOpen(message) {
    const providers = [];

    // Push available providers to an array. Displayed in the dropdown.
    for (const key in this.providersObj.providers) {
      if (key) { providers.push(key); }
    }

    // Open Dropdown Modal
    this.dropdownModalService.open('New Cloud Account', 'Cloud Account', providers);

  }

  openSystemModal(message) {
    this.systemModalService.openSystemModal(message);
  }
  // If the edit button is hit, the Edit modal is opened.
  editCloudAccount() {
    const selectedItems = this.cloudAccountsService.returnSelected();
    let itemindex: string;
    if (selectedItems.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No Cloud Account Selected.');
    } else if (selectedItems.length > 1) {
      this.notifications.display('warn', 'Warning:', 'You cannot edit more than one Cloud Account at a time.');
    } else {
      for (const aprovider in this.providersObj.providers) {
        if (this.providersObj.providers[aprovider]['model']['provider'] === selectedItems[0].provider) {
          this.providersObj.providers[aprovider]['model'] = selectedItems[0];
          itemindex = aprovider;
        }
      }
      this.editModalService.open('Edit', itemindex, this.providersObj.providers);
    }
  }

  // If the delete button is hit, the seleted accounts are deleted.
  deleteCloudAccount() {
    const selectedItems = this.cloudAccountsService.returnSelected();
    if (selectedItems.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No Cloud Account Selected.');
    } else {
      for (const provider of selectedItems) {
        this.supergiant.CloudAccounts.delete(provider.id).subscribe(
          (data) => {
            this.notifications.display('success', 'Cloud Account: ' + provider.name, 'Deleted...');
          },
          (err) => {
            this.notifications.display('error', 'Cloud Account: ' + provider.name, 'Error:' + err);
          },
        );
      }
    }
  }
}

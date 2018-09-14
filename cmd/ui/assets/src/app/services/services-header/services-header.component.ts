import { Component, OnDestroy, AfterViewInit } from '@angular/core';
import { ServicesService } from '../services.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Subscription } from 'rxjs';
import { Notifications } from '../../shared/notifications/notifications.service';
import { SystemModalService } from '../../shared/system-modal/system-modal.service';
import { DropdownModalService } from '../../shared/dropdown-modal/dropdown-modal.service';
import { EditModalService } from '../../shared/edit-modal/edit-modal.service';
import { LoginComponent } from '../../login/login.component';
import { ServicesModel } from '../services.model';

@Component({
  selector: 'app-services-header',
  templateUrl: './services-header.component.html',
  styleUrls: ['./services-header.component.scss']
})
export class ServicesHeaderComponent implements OnDestroy, AfterViewInit {
  providersObj: any;
  subscriptions = new Subscription();
  servicesModel = new ServicesModel;
  kubes = [];
  searchString = '';

  constructor(
    private servicesService: ServicesService,
    private supergiant: Supergiant,
    private notifications: Notifications,
    private systemModalService: SystemModalService,
    private dropdownModalService: DropdownModalService,
    private editModalService: EditModalService,
    public loginComponent: LoginComponent,
  ) { }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  setSearch(value) {
    this.servicesService.searchString = value;
  }
  // After init, grab the schema
  ngAfterViewInit() {
    this.subscriptions.add(this.supergiant.Kubes.get().subscribe(
      (kubes) => { this.kubes = kubes.items; },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); },
    ));

    this.subscriptions.add(this.dropdownModalService.dropdownModalResponse.subscribe(
      (option) => {
        if (option !== 'closed') {
          const kube = this.kubes.filter(resource => resource.name === option)[0];
          this.servicesModel.service.model.kube_name = kube.name;
          this.editModalService.open('Save', 'service', this.servicesModel.providers);
        }
      }, ));


    this.subscriptions.add(this.editModalService.editModalResponse.subscribe(
      (userInput) => {
        if (userInput !== 'closed') {
          const action = userInput[0];
          const providerID = userInput[1];
          const model = userInput[2];
          if (action === 'Edit') {
            this.subscriptions.add(this.supergiant.KubeResources.update(providerID, model).subscribe(
              (data) => { this.success(model); this.servicesService.resetSelected(); },
              (err) => { this.error(model, err); }));
          } else {
            this.subscriptions.add(this.supergiant.KubeResources.create(model).subscribe(
              (data) => { this.success(model); this.servicesService.resetSelected(); },
              (err) => { this.error(model, err); }));
          }
        }
      }));
  }

  success(model) {
    this.notifications.display(
      'success',
      'Service: ' + model.name,
      'Created...',
    );
  }

  error(model, data) {
    this.notifications.display(
      'error',
      'Service: ' + model.name,
      'Error:' + data.statusText);
  }

  sendOpen(message) {
    let options = [];
    options = this.kubes.map((kube) => kube.name);
    this.dropdownModalService.open('New Service', 'Kube', options);

  }

  openSystemModal(message) {
    this.systemModalService.openSystemModal(message);
  }
  // If the edit button is hit, the Edit modal is opened.
  editService() {
    const selectedItems = this.servicesService.returnSelected();

    if (selectedItems.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No Service Selected.');
    } else if (selectedItems.length > 1) {
      this.notifications.display('warn', 'Warning:', 'You cannot edit more than one provider at a time.');
    } else {
      this.providersObj.providers[selectedItems[0].provider].model = selectedItems[0];
      this.editModalService.open('Edit', selectedItems[0].provider, this.providersObj);
    }
  }

  // If the delete button is hit, the seleted accounts are deleted.
  deleteService() {
    const selectedItems = this.servicesService.returnSelected();
    if (selectedItems.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No Service Selected.');
    } else {
      for (const provider of selectedItems) {
        this.subscriptions.add(this.supergiant.KubeResources.delete(provider.id).subscribe(
          (data) => {
            this.notifications.display('success', 'Service: ' + provider.name, 'Deleted...');
            this.servicesService.resetSelected();
          },
          (err) => {
            this.notifications.display('error', 'Service: ' + provider.name, 'Error:' + err);
          },
        ));
      }
    }
  }
}

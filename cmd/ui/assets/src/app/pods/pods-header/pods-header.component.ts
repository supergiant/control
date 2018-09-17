import { Component, OnDestroy, AfterViewInit } from '@angular/core';
import { PodsService } from '../pods.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Subscription } from 'rxjs/Subscription';
import { Notifications } from '../../shared/notifications/notifications.service';
import { SystemModalService } from '../../shared/system-modal/system-modal.service';
import { DropdownModalService } from '../../shared/dropdown-modal/dropdown-modal.service';
import { EditModalService } from '../../shared/edit-modal/edit-modal.service';
import { LoginComponent } from '../../login/login.component';
import { PodsModel } from '../pods.model';

@Component({
  selector: 'app-pods-header',
  templateUrl: './pods-header.component.html',
  styleUrls: ['./pods-header.component.scss']
})
export class PodsHeaderComponent implements OnDestroy, AfterViewInit {
  providersObj: any;
  subscriptions = new Subscription();
  podsModel = new PodsModel;
  kubes = [];
  searchString = '';

  constructor(
    private podsService: PodsService,
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
    this.podsService.searchString = value;
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
          this.podsModel.pod.model.kube_name = kube.name;
          this.editModalService.open('Save', 'pod', this.podsModel.providers);
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
              (data) => {
                this.success(model);
                this.podsService.resetSelected();
              },
              (err) => { this.error(model, err); }));
          } else {
            this.subscriptions.add(this.supergiant.KubeResources.create(model).subscribe(
              (data) => {
                this.success(model);
                this.podsService.resetSelected();
              },
              (err) => { this.error(model, err); }));
          }
        }
      }));
  }

  success(model) {
    this.notifications.display(
      'success',
      'Pod: ' + model.name,
      'Created...',
    );
  }

  error(model, data) {
    this.notifications.display(
      'error',
      'Pod: ' + model.name,
      'Error:' + data.statusText);
  }

  sendOpen(message) {
    let options = [];
    options = this.kubes.map((kube) => kube.name);
    this.dropdownModalService.open('New Pod', 'Kube', options);
  }

  openSystemModal(message) {
    this.systemModalService.openSystemModal(message);
  }
  // If the edit button is hit, the Edit modal is opened.
  editPod() {
    const selectedItems = this.podsService.returnSelected();

    if (selectedItems.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No Pod Selected.');
    } else if (selectedItems.length > 1) {
      this.notifications.display('warn', 'Warning:', 'You cannot edit more than one Pod at a time.');
    } else {
      this.providersObj.providers[selectedItems[0].provider].model = selectedItems[0];
      this.editModalService.open('Edit', selectedItems[0].provider, this.providersObj);
    }
  }

  // If the delete button is hit, the seleted accounts are deleted.
  deletePod() {
    const selectedItems = this.podsService.returnSelected();
    if (selectedItems.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No Pod Selected.');
    } else {
      for (const provider of selectedItems) {
        this.subscriptions.add(this.supergiant.KubeResources.delete(provider.id).subscribe(
          (data) => {
            this.notifications.display('success', 'Pod: ' + provider.name, 'Deleted...');
            this.podsService.resetSelected();
          },
          (err) => {
            this.notifications.display('error', 'Pod: ' + provider.name, 'Error:' + err);
          },
        ));
      }
    }
  }
}

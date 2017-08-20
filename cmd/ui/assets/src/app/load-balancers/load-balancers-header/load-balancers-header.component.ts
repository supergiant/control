import { Component, OnDestroy, AfterViewInit } from '@angular/core';
import { LoadBalancersService } from '../load-balancers.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { LoadBalancersComponent } from '../load-balancers.component';
import { Subscription } from 'rxjs/Subscription';
import { Notifications } from '../../shared/notifications/notifications.service';
import { SystemModalService } from '../../shared/system-modal/system-modal.service';
import { DropdownModalService } from '../../shared/dropdown-modal/dropdown-modal.service';
import { EditModalService } from '../../shared/edit-modal/edit-modal.service';
import { LoginComponent } from '../../login/login.component';
import { LoadBalancersModel } from '../load-balancers.model';

@Component({
  selector: 'app-load-balancers-header',
  templateUrl: './load-balancers-header.component.html',
  styleUrls: ['./load-balancers-header.component.css']
})
export class LoadBalancersHeaderComponent implements OnDestroy, AfterViewInit {
  providersObj: any;
  subscriptions = new Subscription();
  loadBalancersModel = new LoadBalancersModel;
  kubes = [];

  constructor(
    private loadBalancersService: LoadBalancersService,
    private loadBalancersComponent: LoadBalancersComponent,
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
          this.loadBalancersModel.loadBalancer.model.kube_name = kube.name;
          this.editModalService.open('Save', 'loadBalancer', this.loadBalancersModel.providers);
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
            this.subscriptions.add(this.supergiant.LoadBalancers.update(providerID, model).subscribe(
              (data) => {
                this.success(model);
                this.loadBalancersComponent.getAccounts();
              },
              (err) => { this.error(model, err); }));
          } else {
            this.subscriptions.add(this.supergiant.LoadBalancers.create(model).subscribe(
              (data) => {
                this.success(model);
                this.loadBalancersComponent.getAccounts();
              },
              (err) => { this.error(model, err); }));
          }
        }
      }
    ));
  }

  success(model) {
    this.notifications.display(
      'success',
      'Load Balancer: ' + model.name,
      'Created...',
    );
  }

  error(model, data) {
    this.notifications.display(
      'error',
      'Load Balancer: ' + model.name,
      'Error:' + data.statusText);
  }

  // If new button if hit, the New dropdown is triggered.
  sendOpen(message) {
    let options = [];
    options = this.kubes.map((kube) => kube.name);
    this.dropdownModalService.open('New Load Balancer', 'Kube', options);
  }

  openSystemModal(message) {
    this.systemModalService.openSystemModal(message);
  }
  // If the edit button is hit, the Edit modal is opened.
  editUser() {
    const selectedItems = this.loadBalancersService.returnSelected();

    if (selectedItems.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No Load Balancer Selected.');
    } else if (selectedItems.length > 1) {
      this.notifications.display('warn', 'Warning:', 'You cannot edit more than one Load Balancer at a time.');
    } else {
      this.providersObj.providers[selectedItems[0].provider].model = selectedItems[0];
      this.editModalService.open('Edit', selectedItems[0].provider, this.providersObj);
    }
  }

  // If the delete button is hit, the seleted accounts are deleted.
  deleteCloudAccount() {
    const selectedItems = this.loadBalancersService.returnSelected();
    if (selectedItems.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No Load Balancer Selected.');
    } else {
      for (const provider of selectedItems) {
        this.subscriptions.add(this.supergiant.LoadBalancers.delete(provider.id).subscribe(
          (data) => {
            this.notifications.display('success', 'Load Balancer: ' + provider.name, 'Deleted...');
            this.loadBalancersComponent.getAccounts();
          },
          (err) => {
            this.notifications.display('error', 'Load Balancer: ' + provider.name, 'Error:' + err);
          },
        ));
      }
    }
  }
}

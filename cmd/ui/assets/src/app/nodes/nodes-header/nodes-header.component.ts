import { Component, OnDestroy, AfterViewInit } from '@angular/core';
import { NodesService } from '../nodes.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Subscription } from 'rxjs/Subscription';
import { Notifications } from '../../shared/notifications/notifications.service';
import { SystemModalService } from '../../shared/system-modal/system-modal.service';
import { DropdownModalService } from '../../shared/dropdown-modal/dropdown-modal.service';
import { EditModalService } from '../../shared/edit-modal/edit-modal.service';
import { LoginComponent } from '../../login/login.component';
import { NodesModel } from '../nodes.model';

@Component({
  selector: 'app-nodes-header',
  templateUrl: './nodes-header.component.html',
  styleUrls: ['./nodes-header.component.css']
})
export class NodesHeaderComponent implements OnDestroy, AfterViewInit {
  providersObj: any;
  subscriptions = new Subscription();
  nodesModel = new NodesModel;
  kubes = [];
  searchString = '';

  constructor(
    private nodesService: NodesService,
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
    this.nodesService.searchString = value;
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
          this.nodesModel.node.model.kube_name = kube.name;
          for (const size of kube.node_sizes) {
            this.nodesModel.node.schema.properties.size.oneOf = this.nodesModel.node.schema.properties.size.oneOf.concat({
              'enum': [size],
              'description': size,
            }, );
          }
          this.nodesModel.node.schema.properties.size.default = kube.node_sizes[0];
          this.editModalService.open('Save', 'node', this.nodesModel.providers);
        }
      }, ));

    this.subscriptions.add(this.editModalService.editModalResponse.subscribe(
      (userInput) => {
        if (userInput !== 'closed') {
          const action = userInput[0];
          const providerID = userInput[1];
          const model = userInput[2];
          if (action === 'Edit') {
            this.subscriptions.add(this.supergiant.Nodes.update(providerID, model).subscribe(
              (data) => {
                this.success(model);
                this.nodesService.resetSelected();
              },
              (err) => { this.error(model, err); }));
          } else {
            this.subscriptions.add(this.supergiant.Nodes.create(model).subscribe(
              (data) => {
                this.success(model);
                this.nodesService.resetSelected();
              },
              (err) => { this.error(model, err); }));
          }
        }
      }));
  }

  success(model) {
    this.notifications.display(
      'success',
      'Node: ' + model.provider_id,
      'Created...',
    );
  }

  error(model, data) {
    this.notifications.display(
      'error',
      'Node: ' + model.provider_id,
      'Error:' + data.statusText);
  }

  sendOpen(message) {
    let options = [];
    options = this.kubes.map((kube) => kube.name);
    this.dropdownModalService.open('New Node', 'Kube', options);
  }

  openSystemModal(message) {
    this.systemModalService.openSystemModal(message);
  }
  // If the edit button is hit, the Edit modal is opened.
  editNode() {
    const selectedItems = this.nodesService.returnSelected();

    if (selectedItems.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No Node Selected.');
    } else if (selectedItems.length > 1) {
      this.notifications.display('warn', 'Warning:', 'You cannot edit more than one Node at a time.');
    } else {
      this.providersObj.providers[selectedItems[0].provider].model = selectedItems[0];
      this.editModalService.open('Edit', selectedItems[0].provider, this.providersObj);
    }
  }

  // If the delete button is hit, the seleted accounts are deleted.
  deleteNode() {
    const selectedItems = this.nodesService.returnSelected();
    if (selectedItems.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No Node Selected.');
    } else {
      for (const node of selectedItems) {
        this.subscriptions.add(this.supergiant.Nodes.delete(node.id).subscribe(
          (data) => {
            this.notifications.display('success', 'Node: ' + node.provider_id, 'Deleted...');
            this.nodesService.resetSelected();
          },
          (err) => {
            this.notifications.display('error', 'Node: ' + node.provider_id, 'Error:' + err);
          },
        ));
      }
    }
  }
}

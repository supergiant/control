import { Component, OnDestroy, AfterViewInit } from '@angular/core';
import { AppsService } from '../apps.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { AppsComponent } from '../apps.component';
import { Subscription } from 'rxjs/Subscription';
import { Notifications } from '../../shared/notifications/notifications.service';
import { SystemModalService } from '../../shared/system-modal/system-modal.service';
import { DropdownModalService } from '../../shared/dropdown-modal/dropdown-modal.service';
import { EditModalService } from '../../shared/edit-modal/edit-modal.service';
import { LoginComponent } from '../../login/login.component';
import { RepoModalService } from '../repo-modal/repo-modal.service';
import { AppsModel } from '../apps.model';
import * as GenerateSchema from 'generate-schema';


@Component({
  selector: 'app-apps-header',
  templateUrl: './apps-header.component.html',
  styleUrls: ['./apps-header.component.css']
})
export class AppsHeaderComponent implements OnDestroy, AfterViewInit {
  subscriptions = new Subscription();
  kubes = [];
  appsModel = new AppsModel;

  constructor(
    private appsService: AppsService,
    private appsComponent: AppsComponent,
    private supergiant: Supergiant,
    private notifications: Notifications,
    private systemModalService: SystemModalService,
    private repoModalService: RepoModalService,
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
          const chart = this.appsService.returnSelected();
          if (chart.length === 0) {
            this.notifications.display('warn', 'Warning:', 'No App Selected.');
          } else if (chart.length > 1) {
            this.notifications.display('warn', 'Warning:', 'You cannot deploy more than on App at a time.');
          } else {
            if (chart[0].default_config) {
              // this is our model: the vars file provided by the chart.
              this.appsModel.app.model.config = chart[0].default_config;
              this.appsModel.app.model.chart_name = chart[0].name;
              this.appsModel.app.model.chart_version = chart[0].version;
              this.appsModel.app.model.kube_name = kube.name;
              this.appsModel.app.model.repo_name = chart[0].repo_name;
              // We dynamically generate a schema from the vars file.
              // TODO: Note - we need to add a sniffer here to look for a schema file
              // and any icon images in the chart. If they exist, we should use them instead of the generated one.
              this.appsModel.app.schema = GenerateSchema.json(this.appsModel.app.model);
              this.editModalService.open('Save', 'app', this.appsModel);
            } else {
              // Some charts do not have var files. In that case there is no need for a editor at all...
              // In this case we skip the editor and just deploy.
              this.appsModel.app.model.chart_name = chart[0].name;
              this.appsModel.app.model.chart_version = chart[0].version;
              this.appsModel.app.model.kube_name = kube.name;
              this.appsModel.app.model.repo_name = chart[0].repo_name;
              this.editModalService.editModalResponse.next(['Save', 'app', this.appsModel.app.model]);
              console.log('Launch config... no edit box needed');
            }
          }
        }
      }
    ));

    this.subscriptions.add(this.editModalService.editModalResponse.subscribe(
      (userInput) => {
        if (userInput !== 'closed') {
          const action = userInput[0];
          const providerID = userInput[1];
          const model = userInput[2];
          if (action === 'Edit') {
            this.subscriptions.add(this.supergiant.HelmReleases.update(providerID, model).subscribe(
              (data) => {
                this.success(model);
                this.appsComponent.getDeployments();
              },
              (err) => { this.error(model, err); }));
          } else {
            this.subscriptions.add(this.supergiant.HelmReleases.create(model).subscribe(
              (data) => {
                this.success(model);
                this.appsComponent.getDeployments();
              },
              (err) => { this.error(model, err); }));
          }
        }
      }));
  }

  success(model) {
    this.notifications.display(
      'success',
      'App: ' + model.chart_name,
      'Created...',
    );
  }

  error(model, data) {
    this.notifications.display(
      'error',
      'App: ' + model.chart_name,
      'Error:' + data.statusText);
  }
  // If new button if hit, the New dropdown is triggered.
  sendOpen(message) {
    let options = [];
    options = this.kubes.map((kube) => kube.name);
    this.dropdownModalService.open('New App', 'Kubes', options);
  }

  openSystemModal(message) {
    this.systemModalService.openSystemModal(message);
  }

  openRepoModal(message) {
    this.repoModalService.openRepoModal(message);
  }
  // If the edit button is hit, the Edit modal is opened.
  editKube() {
    // var selectedItems = this.appsService.returnSelected()
    //
    // if (selectedItems.length === 0) {
    //   this.notifications.display("warn", "Warning:", "No Provider Selected.")
    // } else if (selectedItems.length > 1) {
    //   this.notifications.display("warn", "Warning:", "You cannot edit more than one provider at a time.")
    // } else {
    //   this.providersObj.providers[selectedItems[0].provider].model = selectedItems[0]
    //   this.editModalService.open("Edit", selectedItems[0].provider, this.providersObj);
    // }
  }

  // If the delete button is hit, the seleted accounts are deleted.
  deleteCloudAccount() {
    const selectedItems = this.appsService.returnSelected();
    if (selectedItems.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No App Selected.');
    } else {
      for (const provider of selectedItems) {
        this.subscriptions.add(this.supergiant.HelmReleases.delete(provider.id).subscribe(
          (data) => {
            this.notifications.display('success', 'App: ' + provider.name, 'Deleted...');
          },
          (err) => {
            this.notifications.display('error', 'App: ' + provider.name, 'Error:' + err);
          },
        ));
      }
    }
  }
}

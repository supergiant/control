import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { KubeResourcesModel } from '../kube-resources.model';

import "brace/mode/json";

@Component({
  selector: 'app-new-kube-resource',
  templateUrl: './new-kube-resource.component.html',
  styleUrls: ['./new-kube-resource.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class NewKubeResourceComponent implements OnInit {

  subscriptions = new Subscription();
  kubeResourcesModel = new KubeResourcesModel();
  resourceTypes = [
      { displayName: 'Pod', type: 'pod' },
      { displayName: 'Service', type: 'service' },
      { displayName: 'LoadBalancer', type: 'loadBalancer' },
    ];
  private selectedResourceType: string;
  private kubeId: number;
  private kubeName: string;
  private textStatus = 'form-control';
  private schema: any;
  private model: any;
  private value: any;
  private modelString: string;
  private layout: any;
  private badString: string;
  private isDisabled: boolean;

  constructor(
    private supergiant: Supergiant,
    private notifications: Notifications,
    private route: ActivatedRoute,
    private router: Router,
  ) { }

  updateModelKubeName() {
    // TODO: burn this with fire and find another way
    this.model.kube_name = this.kubeName;
    let ms = JSON.parse(this.modelString);
    ms.kube_name = this.kubeName;
    this.modelString = JSON.stringify(ms, null, 2);
  }

  getClusterName(kubeId) {
    this.subscriptions.add(this.supergiant.Kubes.get(kubeId).subscribe(
      (kube) => {
        this.kubeName = kube.name;
        this.updateModelKubeName();
      }
    ))
  }

  createKubeResource(model) {
    // some attr keys are dynamic, which angular2-json-schema-form doesn't support _eyeroll_
    let resourceService: any;
    switch (model.kind) {
      case 'Pod': {
        if (model.resource.metadata) {
          model.resource.metadata.labels = this.parseArrToObj(model.resource.metadata.labels);
        }
        resourceService = this.supergiant.KubeResources;
        break;
      }
      case 'Service': {
        model.template.spec.selector = this.parseArrToObj(model.template.spec.selector);
        resourceService = this.supergiant.KubeResources;
        break;
      }
      case 'LoadBalancer': {
        model.selector = this.parseArrToObj(model.selector);
        model.ports = this.parseArrToObj(model.ports);
        resourceService = this.supergiant.LoadBalancers;
        break;
      }
    }
    this.subscriptions.add(resourceService.create(model).subscribe(
      (data) => {
        this.success(model);
        this.router.navigate(['/clusters/', this.kubeId]);
      },
      (err) => { this.error(model, err); }));

  }

  success(model) {
    this.notifications.display(
      'success',
      'Resource: ' + model.name,
      'Created...',
    );
  }

  error(model, data) {
    this.notifications.display(
      'error',
      'Resource: ' + model.name,
      'Error:' + data.statusText);
  }

  parseArrToObj(arr) {
    if (arr) {
      let returnObj = arr.reduce(
        (obj, label) => {
          obj[label['key']] = label['value'];
          return obj;
        },
        {});
      return returnObj;
    } else { return {} }
  }

  updateFromForm(model) {
    this.model = model;
    this.modelString = JSON.stringify(model, null, 2);
  }

  convertToObj(json) {
    try {
      JSON.parse(json);
    } catch (e) {
      this.textStatus = 'form-control badTextarea';
      this.badString = e;
      this.isDisabled = true;
      return;
    }
    this.textStatus = 'form-control goodTextarea';
    this.badString = 'Valid JSON';
    this.isDisabled = false;
    this.model = JSON.parse(json);
    this.modelString = JSON.stringify(this.model, null, 2);
  }

  resetModel(selectedResource) {
    switch (selectedResource) {
      case "Pod": {
        this.chooseResourceType({ displayName: 'Pod', type: 'pod' });
        break;
      }
      case "Service": {
        this.chooseResourceType({ displayName: 'Service', type: 'service' });
        break;
      }
      case "LoadBalancer": {
        this.chooseResourceType({ displayName: 'LoadBalancer', type: 'loadBalancer' });
        break;
      }
    }
  }

  reset() {
    this.model = null;
    this.modelString = null;
    this.schema = null;
    this.layout = null;
  }

  back() {
    window.history.back();
  }

  chooseResourceType(resource) {
    this.schema = this.kubeResourcesModel[resource.type].schema;
    this.layout = this.kubeResourcesModel[resource.type].layout;
    this.model = this.kubeResourcesModel[resource.type].model;
    this.modelString = JSON.stringify(this.model, null, 2);
    this.selectedResourceType = resource.displayName;
    this.getClusterName(this.kubeId);
  }

  ngOnInit() {
    this.kubeId = this.route.snapshot.params.id;
  }

}

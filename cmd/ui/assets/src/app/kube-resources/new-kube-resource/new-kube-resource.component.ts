import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { KubeResourcesModel } from '../kube-resources.model';

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
  private textStatus = 'form-control';
  private schema: any;
  private model: any;
  private layout: any;
  private badString: string;
  private isDisabled: boolean;

  constructor(
    private supergiant: Supergiant,
    private notifications: Notifications,
    private route: ActivatedRoute,
    private router: Router,
  ) { }

  getClusterName(kubeId) {
    this.subscriptions.add(this.supergiant.Kubes.get(kubeId).subscribe(
      (kube) => {
        this.model.kube_name = kube.name;
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
  }

  reset() {
    this.model = null;
    this.schema = null;
    this.layout = null;
  }

  back() {
    window.history.back();
  }

  chooseResourceType(resource) {
    switch (resource.type) {
      case 'pod': {
        this.schema = this.kubeResourcesModel.pod.schema;
        this.layout = this.kubeResourcesModel.pod.layout;
        this.model = this.kubeResourcesModel.pod.model;
        this.selectedResourceType = resource.displayName;
        this.getClusterName(this.kubeId)
        break;
      }
      case 'service': {
        this.schema = this.kubeResourcesModel.service.schema;
        this.layout = this.kubeResourcesModel.service.layout;
        this.model = this.kubeResourcesModel.service.model;
        this.selectedResourceType = resource.displayName;
        this.getClusterName(this.kubeId)
        break;
      }
      case 'loadBalancer': {
        this.schema = this.kubeResourcesModel.loadBalancer.schema;
        this.layout = this.kubeResourcesModel.loadBalancer.layout;
        this.model = this.kubeResourcesModel.loadBalancer.model;
        this.selectedResourceType = resource.displayName;
        this.getClusterName(this.kubeId)
        break;
      }
      default: {
        this.schema = null;
        this.layout = null;
        this.model = null;
        break;
      }
    }


  }

  ngOnInit() {
    this.kubeId = this.route.snapshot.params.id;
  }

}

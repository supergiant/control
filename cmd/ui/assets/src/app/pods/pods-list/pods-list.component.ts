
import {timer as observableTimer,  Subscription ,  Observable } from 'rxjs';

import {switchMap} from 'rxjs/operators';
import { Component, OnInit, Input, OnDestroy, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { ChartsModule, BaseChartDirective } from 'ng2-charts';
import { ContextMenuService, ContextMenuComponent } from 'ngx-contextmenu';
import { PodsModel } from '../pods.model';


@Component({
  selector: 'app-pods-list',
  templateUrl: './pods-list.component.html',
  styleUrls: ['./pods-list.component.scss']
})
export class PodsListComponent implements OnInit, OnDestroy {

  @Input() kube: any;
  @ViewChild(ContextMenuComponent) public basicMenu: ContextMenuComponent;

  private noderows = [];
  private newNode = false;
  public nodecolumns: Array<any> = [];
  private subscriptions = new Subscription();
  public resources = [];

  // TODO: Kinds should be dynamically built.
  public resouceKinds = ['Pod', 'Service', 'PersistentVolume', 'LoadBalancer'];
  public selectedResourceKind = 'Pod';
  private rowChartLegend = false;
  private rowChartType = 'line';
  private rowChartLabels: Array<any> = ['', '', '', '', '', '', ''];
  public selected = [];
  private rawEvent: any;
  private contextmenuRow: any;
  private contextmenuColumn: any;
  private selectedSize: string;
  private nodeName: string;
  private podsModel = new PodsModel;
  public displayCheck: boolean;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private supergiant: Supergiant,
    private notifications: Notifications,
    private contextMenuService: ContextMenuService,
  ) { }

  ngOnInit() {
    this.get();
  }

  success(model) {
    this.notifications.display(
      'success',
      'Node: ' + model.kube_name,
      'Created...',
    );
  }

  error(model, data) {
    this.notifications.display(
      'error',
      'Node: ' + model.kube_name,
      'Error:' + data.statusText);
  }

  get() {
    this.subscriptions.add(observableTimer(0, 5000).pipe(
      switchMap(() => this.supergiant.KubeResources.get())).subscribe(
        (pods) => {

          // Filter by kube if it is passed, else no filter.
          if (this.kube) {
            this.resources = pods.items.filter(
              resource =>
                resource.kube_name === this.kube.name
                && resource.kind === this.selectedResourceKind
                || resource.kube_name === this.kube.name
                && resource.resource.spec.type === this.selectedResourceKind
            );
          } else {
            this.resources = pods.items.filter(
              resource =>
                resource.kind === this.selectedResourceKind
                || resource.resource.spec.type === this.selectedResourceKind
            );
          }


          this.resources.map(resource => ({
            id: resource.id,
            name: resource.name,
            namespace: resource.namespace,
            status: resource.passive_status,
          })
          );

          const selected: Array<any> = [];
          this.selected.forEach((kube, index) => {
            for (const row of this.resources) {
              if (row.id === kube.id) {
                selected.push(row);
                break;
              }
            }
          });
          this.selected = selected;
        },
        (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));

  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  onActivate(activated) {

    if (activated.type === 'click' && activated.column.name !== 'checkbox') {
      this.router.navigate(['/pods', activated.row.id]);
    }
  }

  onSelect({ selected }) {
    this.selected.splice(0, this.selected.length);
    this.selected.push(...selected);
  }

  usageOrZeroCPU(usage) {
    if (usage == null) {
      return ([0, 0, 0, 0, 0, 0, 0, 0, 0, 0]);
    } else {
      return usage.cpu_usage_rate.map((data) => data.value);
    }
  }

  onTableContextMenu(contextMenuEvent) {
    this.rawEvent = contextMenuEvent.event;
    if (contextMenuEvent.type === 'body') {
      this.contextmenuColumn = undefined;
      this.contextMenuService.show.next({
        contextMenu: this.basicMenu,
        item: contextMenuEvent.content,
        event: contextMenuEvent.event,
      });
    } else {
      this.contextmenuColumn = contextMenuEvent.content;
      this.contextmenuRow = undefined;
    }

    contextMenuEvent.event.preventDefault();
    contextMenuEvent.event.stopPropagation();
  }

  contextDelete(item) {
    for (const row of this.resources) {
      if (row.id === item.id) {
        this.selected.push(row);
        this.delete();
        break;
      }
    }
  }

  delete() {
    if (this.selected.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No Pod Selected.');
    } else {
      for (const pod of this.selected) {
        this.subscriptions.add(this.supergiant.KubeResources.delete(pod.id).subscribe(
          (data) => {
            this.notifications.display('success', 'Pod: ' + pod.name, 'Deleted...');
            this.selected = [];
          },
          (err) => {
            this.notifications.display('error', 'Pods: ' + pod.name, 'Error:' + err);
          },
        ));
      }
    }
  }
}

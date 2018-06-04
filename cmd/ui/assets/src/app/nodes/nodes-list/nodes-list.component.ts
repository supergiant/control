import { Component, OnInit, Input, OnDestroy, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { Observable } from 'rxjs/Observable';
import { ChartsModule, BaseChartDirective } from 'ng2-charts';
import { ContextMenuService, ContextMenuComponent } from 'ngx-contextmenu';
import { NodesModel } from './nodes.model';

@Component({
  selector: 'app-nodes-list',
  templateUrl: './nodes-list.component.html',
  styleUrls: ['./nodes-list.component.scss']
})
export class NodesListComponent implements OnInit, OnDestroy {

  @Input() kube: any;
  @ViewChild(ContextMenuComponent) public basicMenu: ContextMenuComponent;

  private noderows = [];
  private newNode = false;
  public nodecolumns: Array<any> = [];
  private subscriptions = new Subscription();
  public nodes = [];

  // Usage chart setup:
  private rowChartLegend = false;
  private rowChartType = 'line';
  private rowChartLabels: Array<any> = ['', '', '', '', '', '', ''];
  private rowChartOptions: any = {
    responsive: true,
    maintainAspectRatio: true,
  };
  private rowChartColors: Array<any> = [
    { // grey
      backgroundColor: 'rgba(148,159,177,0.2)',
      borderColor: 'rgba(148,159,177,1)',
      pointBackgroundColor: 'rgba(148,159,177,1)',
      pointBorderColor: '#fff',
      pointHoverBackgroundColor: '#fff',
      pointHoverBorderColor: 'rgba(148,159,177,0.8)'
    },
    { // dark grey
      backgroundColor: 'rgba(77,83,96,0.2)',
      borderColor: 'rgba(77,83,96,1)',
      pointBackgroundColor: 'rgba(77,83,96,1)',
      pointBorderColor: '#fff',
      pointHoverBackgroundColor: '#fff',
      pointHoverBorderColor: 'rgba(77,83,96,1)'
    },
    { // grey
      backgroundColor: 'rgba(148,159,177,0.2)',
      borderColor: 'rgba(148,159,177,1)',
      pointBackgroundColor: 'rgba(148,159,177,1)',
      pointBorderColor: '#fff',
      pointHoverBackgroundColor: '#fff',
      pointHoverBorderColor: 'rgba(148,159,177,0.8)'
    }
  ];

  public selected = [];
  private rawEvent: any;
  private contextmenuRow: any;
  private contextmenuColumn: any;
  private selectedSize: string;
  private nodeName: string;
  private nodesModel = new NodesModel;
  public displayCheck: boolean;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private supergiant: Supergiant,
    private notifications: Notifications,
    private contextMenuService: ContextMenuService,
  ) { }

  ngOnInit() {
    this.getNodes();
  }

  createNode() {
    const model = this.nodesModel.node.model;
    model.kube_name = this.kube.name;
    model.size = this.selectedSize;

    this.subscriptions.add(this.supergiant.Nodes.create(model).subscribe(
      (data) => {
        this.success(model);
      },
      (err) => { this.error(model, err); }));
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

  getNodes() {
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.Nodes.get()).subscribe(
        (nodes) => {
          this.nodes = nodes.items.filter(
            node => node.kube_name === this.kube.name
          ).map(node => ({
            id: node.id,
            name: node.name,
            size: node.size,
            ip: node.external_ip,
            chartData: [
              { label: 'CPU Usage', data: this.usageOrZeroCPU(node.extra_data) },
            ],
          })
          );

          const selected: Array<any> = [];
          this.selected.forEach((kube, index) => {
            for (const row of this.nodes) {
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
      this.router.navigate(['/nodes', activated.row.id]);
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
    for (const row of this.nodes) {
      if (row.id === item.id) {
        this.selected.push(row);
        this.delete();
        break;
      }
    }
  }

  delete() {
    if (this.selected.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No Node Selected.');
    } else {
      for (const node of this.selected) {
        this.subscriptions.add(this.supergiant.Nodes.delete(node.id).subscribe(
          (data) => {
            this.notifications.display('success', 'Node: ' + node.name, 'Deleted...');
            this.selected = [];
          },
          (err) => {
            this.notifications.display('error', 'Node: ' + node.name, 'Error:' + err);
          },
        ));
      }
    }
  }

}

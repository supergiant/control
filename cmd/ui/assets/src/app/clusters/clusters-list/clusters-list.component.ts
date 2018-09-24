
import {timer as observableTimer,  Subscription ,  Observable } from 'rxjs';

import {switchMap} from 'rxjs/operators';
import { Component, OnInit, OnDestroy, TemplateRef, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { TitleCasePipe } from '@angular/common';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';

import { ContextMenuService, ContextMenuComponent } from 'ngx-contextmenu';

@Component({
  selector: 'app-clusters-list',
  templateUrl: './clusters-list.component.html',
  styleUrls: ['./clusters-list.component.scss'],

})
export class ClustersListComponent implements OnInit, OnDestroy {
  @ViewChild(ContextMenuComponent) public basicMenu: ContextMenuComponent;
  hasCluster = false;
  hasCloudAccount = false;
  hasApp = false;
  clusterCount = 0;
  appCount = 0;
  filterText = '';
  unfilteredRows: Array<any> = [];
  rows: Array<any> = [];
  selected: Array<any> = [];
  columns: Array<any> = [];

  public rowChartOptions: any = {
    responsive: true,
    maintainAspectRatio: true,
    scales: {
      yAxes: [{
        ticks: {
          beginAtZero: true,
        }
      }]
    }
  };
  public rowChartColors: Array<any> = [
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

  // linter is angry about the boolean typing but without it charts
  public rowChartLegend = false;
  public rowChartType = 'line';
  public rowChartLabels: Array<any> = ['', '', '', '', '', '', ''];

  private subscriptions = new Subscription();
  public kubes = [];
  rawEvent: any;
  contextmenuRow: any;
  contextmenuColumn: any;
  constructor(
    private supergiant: Supergiant,
    private notifications: Notifications,
    private titleCase: TitleCasePipe,
    private contextMenuService: ContextMenuService,
    private router: Router,
  ) { }

  ngOnInit() {
    this.getKubes();
    this.getCloudAccounts();
    this.getClusters();
    this.getDeployments();
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  filterRows(filterRows: Array<any>, filterText: string): Array<any> {
    if (filterText === '') {
      return filterRows;
    }
    const matchingRows = [];
    for (const row of filterRows) {
      for (const key of Object.keys(row)) {
        const value = row[key].toString().toLowerCase();
        if (value.toString().indexOf(filterText.toLowerCase()) >= 0) {
          matchingRows.push(row);
          break;
        }
      }
    }
    return matchingRows;
  }

  keyUpFilter(filterText) {
    this.filterText = filterText;
    this.rows = this.filterRows(this.unfilteredRows, filterText);
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

  // public onContextMenu($event: MouseEvent, item: any): void {
  //     this.contextMenuService.show.next({
  //       // Optional - if unspecified, all context menu components will open
  //       contextMenu: this.contextMenu,
  //       event: $event,
  //       item: item,
  //     });
  //     $event.preventDefault();
  //     $event.stopPropagation();
  //   }

  onSelect({ selected }) {
    this.selected.splice(0, this.selected.length);
    this.selected.push(...selected);
  }

  onActivate(activated) {
    if (activated.type === 'click' && activated.column.name !== 'checkbox') {
      this.router.navigate(['/clusters', activated.row.id]);
    }
  }

  lengthOrZero(lenobj) {
    if (lenobj == null) {
      return 0;
    } else {
      return Object.keys(lenobj).length;
    }
  }

  progressOrDone(progobj) {
    if (progobj.status == null) {
      return 'Running';
    } else {
      return progobj.status.description;
    }
  }

  usageOrZeroCPU(usage) {
    if (usage == null) {
      return ([0, 0, 0, 0, 0, 0, 0, 0, 0, 0]);
    } else {
      return usage.cpu_usage_rate.map((data) => data.value);
    }
  }

  getCloudAccounts() {
    this.subscriptions.add(this.supergiant.CloudAccounts.get().subscribe(
      (cloudAccounts) => {
        if (Object.keys(cloudAccounts).length > 0) { this.hasCloudAccount = true; }
      })
    );
  }

  getClusters() {
    this.subscriptions.add(this.supergiant.Kubes.get().subscribe(
      (clusters) => {
        if (Object.keys(clusters.items).length > 0) {
          this.hasCluster = true;
          // this.lineChartData[0]['data'].length = 0;
          // this.lineChartData[0]['data'].length = 0;
          for (const cluster of clusters.items) {
            // this.getKubes(cluster.id);
            this.getKubes();
          }
          this.clusterCount = Object.keys(clusters.items).length;
        }
      })
    );
  }

  getDeployments() {
    this.subscriptions.add(this.supergiant.HelmReleases.get().subscribe(
      (deployments) => {
        if (Object.keys(deployments.items).length > 0) {
          this.hasApp = true;
          this.appCount = Object.keys(deployments.items).length;
        }
      })
    );
  }

  getKubes() {
    this.subscriptions.add(observableTimer(0, 5000).pipe(
      switchMap(() => this.supergiant.Kubes.get())).subscribe(
        (kubes) => {

          const rows = kubes.items.map(kube => ({
            id: kube.id,
            name: kube.name,
            version: kube.kubernetes_version,
            cloudaccount: kube.cloud_account_name,
            nodes: this.lengthOrZero(kube.nodes),
            apps: this.lengthOrZero(kube.helm_releases),
            status: this.titleCase.transform(this.progressOrDone(kube)),
            kube: kube,
            chartData: [
              { label: 'CPU Usage', data: this.usageOrZeroCPU(kube.extra_data) },
              // this should be set to the length of largest array.
            ],
          }));
          // Copy over any kubes that happen to be currently selected.
          const selected = [];
          this.selected.forEach((kube, index) => {
            for (const row of rows) {
              if (row.id === kube.id) {
                selected.push(row);
                break;
              }
            }
          });
          this.unfilteredRows = rows;
          if (Object.keys(rows).length > 0) { this.hasCluster = true; }
          this.rows = this.filterRows(rows, this.filterText);
          this.selected = selected;
        },
        (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  contextDelete(item) {
    for (const row of this.rows) {
      if (row.id === item.id) {
        this.selected.push(row);
        this.deleteKube();
        break;
      }
    }
  }

  deleteKube() {
    if (this.selected.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No Kube Selected.');
    } else {
      for (const provider of this.selected) {
        this.subscriptions.add(this.supergiant.Kubes.delete(provider.id).subscribe(
          (data) => {
            this.notifications.display('success', 'Kube: ' + provider.name, 'Deleted...');
            this.selected = [];
          },
          (err) => {
            this.notifications.display('error', 'Kube: ' + provider.name, 'Error:' + err);
          },
        ));
      }
    }
  }

}

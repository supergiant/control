import { Component, OnInit, OnDestroy, TemplateRef, ViewChild } from '@angular/core';
import { TitleCasePipe } from '@angular/common';
import { Subscription } from 'rxjs/Subscription';
import { Observable } from 'rxjs/Observable';
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
  rows = [];
  selected = [];
  columns = [];

  public rowChartOptions: any = {
    responsive: false
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


  public rowChartLegend: boolean = false;
  public rowChartType: string = 'line';
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
  ) { }

  ngOnInit() {
    this.getKubes();
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  onTableContextMenu(contextMenuEvent) {
      this.rawEvent = contextMenuEvent.event;
      if (contextMenuEvent.type === 'body') {
        console.log(contextMenuEvent);
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
      return( [0, 0, 0, 0, 0, 0, 0, 0, 0, 0] );
    } else {
      return usage.cpu_usage_rate.map((data) => data.value);
    }
  }

  usageOrZeroMEM(usage) {
    if (usage == null) {
      return( [0, 0, 0, 0, 0, 0, 0, 0, 0, 0] );
    } else {
      return usage.memory_usage.map((data) => data.value / 1073741824);
    }
  }

  getKubes() {
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.Kubes.get()).subscribe(
      (kubes) => {

        const rows = kubes.items.map(kube => ({
          id: kube.id,
          name: kube.name,
          version: kube.kubernetes_version,
          cloudaccount: kube.cloud_account_name,
          nodes: this.lengthOrZero(kube.nodes),
          apps: this.lengthOrZero(kube.helmreleases),
          status: this.titleCase.transform(this.progressOrDone(kube)),
          kube: kube,
          chartData: [
            { label: 'CPU Usage', data: this.usageOrZeroCPU(kube.extra_data) },
            { label: 'RAM Usage', data: this.usageOrZeroMEM(kube.extra_data) },
            // this should be set to the length of largest array.
          ],
        }));
        console.log(rows);
        // Copy over any kubes that happen to be currently selected.
        const selected: Array<any> = [];
        this.selected.forEach((kube, index) => {
          for (const row of rows) {
            if (row.id === kube.id) {
              selected.push(row);
              break;
            }
          }
        });
        this.rows = rows;
        this.selected = selected;
      },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  contextDelete(item) {
    console.log(item);
    for (const row of this.rows) {
      if (row.id === item.id) {
        this.selected.push(row);
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

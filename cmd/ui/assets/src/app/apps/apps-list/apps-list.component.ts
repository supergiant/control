import { Component, OnInit, OnDestroy, Pipe, PipeTransform, TemplateRef, ViewChild, Input } from '@angular/core';
// import { Observable } from 'rxjs/Observable';
import { timer } from 'rxjs/observable/timer';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { AppsService } from '../apps.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { ContextMenuService, ContextMenuComponent } from 'ngx-contextmenu';

@Component({
  selector: 'app-apps-list',
  templateUrl: './apps-list.component.html',
  styleUrls: ['./apps-list.component.scss']
})
export class AppsListComponent implements OnInit, OnDestroy {
  @ViewChild(ContextMenuComponent) public basicMenu: ContextMenuComponent;
  @Input() kube: any;
  public selected: Array<any> = [];
  public rows: Array<any> = [];
  public columns: Array<any> = [];
  public displayCheck: boolean;
  private subscriptions = new Subscription();
  public unfilteredRows: Array<any> = [];
  public filterText = '';
  private rawEvent: any;
  contextmenuRow: any;
  contextmenuColumn: any;
  constructor(
    private supergiant: Supergiant,
    private notifications: Notifications,
    private contextMenuService: ContextMenuService,
  ) { }

  ngOnInit() {
    this.getApps();
  }

  // getCharts() {
  //   this.subscriptions.add(Observable.timer(0, 5000)
  //     .switchMap(() => this.supergiant.HelmCharts.get()).subscribe(
  //     (apps) => { this.apps = apps.items; },
  //     () => { }));
  // }

  getApps() {

    this.subscriptions.add(timer(0, 10000)
      .switchMap(() => this.supergiant.HelmReleases.get()).subscribe(
        (deployments) => {
          if (this.kube) {
            this.rows = deployments.items.filter(
              deployment =>
                deployment.kube_name === this.kube.name
            );
          } else {
            this.rows = deployments.items;
          }

          this.rows.map(deployment => ({
            id: deployment.id,
            name: deployment.name,
            kube_name: deployment.kube_name,
            revision: deployment.revision,
            chart_name: deployment.chart_name,
            chart_version: deployment.chart_version,
            updated_value: deployment.updated_value,
            status_value: deployment.status_value
          })
          );

          const selected: Array<any> = [];
          this.selected.forEach((kube, index) => {
            for (const row of this.rows) {
              if (row.id === kube.id) {
                selected.push(row);
                break;
              }
            }
          });
          this.selected = selected;
        },
        () => { }));
  }

  deleteApp() {
    if (this.selected.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No App Selected.');
    } else {
      for (const provider of this.selected) {
        this.subscriptions.add(this.supergiant.HelmReleases.delete(provider.id).subscribe(
          (data) => {
            this.notifications.display('success', 'App: ' + provider.name, 'Deleted...');
            this.selected = [];
          },
          (err) => {
            this.notifications.display('error', 'App: ' + provider.name, 'Error:' + err);
          },
        ));
      }
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

  onActivate(activated) {
    if (activated.type === 'click') {
      // this.router.navigate(['/apps', activated.row.id]);
    }
  }

  filterRows(filterRows: Array<any>, filterText: string): Array<any> {
    if (filterText === '') {
      return filterRows;
    }
    const matchingRows = [];
    for (const row of filterRows) {
      for (const key of Object.keys(row)) {
        if (row[key] != null) {
          const value = row[key].toString().toLowerCase();
          if (value.toString().indexOf(filterText.toLowerCase()) >= 0) {
            matchingRows.push(row);
            break;
          }
        }
      }
    }
    return matchingRows;
  }

  keyUpFilter(filterText) {
    this.filterText = filterText;
    this.rows = this.filterRows(this.unfilteredRows, filterText);
  }

  contextDelete(item) {
    for (const row of this.rows) {
      if (row.id === item.id) {
        this.selected.push(row);
        this.deleteApp();
        break;
      }
    }
  }

  onSelect({ selected }) {
    this.selected.splice(0, this.selected.length);
    this.selected.push(...selected);
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }
}

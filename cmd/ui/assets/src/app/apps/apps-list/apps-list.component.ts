import { Component, OnInit, OnDestroy, Pipe, PipeTransform, TemplateRef, ViewChild } from '@angular/core';
import { Observable } from 'rxjs/Observable';
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
  public selected: Array<any> = [];
  public rows: Array<any> = [];
  private subscriptions = new Subscription();
  public unfilteredRows: Array<any> = [];
  public filterText: string = '';
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
    this.subscriptions.add(Observable.timer(0, 10000)
      .switchMap(() => this.supergiant.HelmReleases.get()).subscribe(
      (deployments) => {
        const selected: Array<any> = [];
        this.selected.forEach((app, index) => {
          for (const row of deployments.items) {
            if (row.id === app.id) {
              selected.push(row);
              break;
            }
          }
        });
        this.unfilteredRows = deployments.items;
        this.rows = this.filterRows(deployments.items, this.filterText);
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

  onActivate(activated) {
    if (activated.type === 'click') {
      // this.router.navigate(['/apps', activated.row.id]);
    }
  }

  filterRows(filterRows: Array<any>, filterText: string): Array<any> {
    console.log(filterRows);
    console.log(filterText);
    if (filterText === '') {
      return filterRows;
    }
    const matchingRows = [];
    for (const row of filterRows) {
      for (const key of Object.keys(row)) {
        if ( row[key] != null) {
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
    console.log(item);
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

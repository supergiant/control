import { Component, OnInit, OnDestroy, ViewEncapsulation, ViewChild } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { Observable } from 'rxjs/Observable';
import { Supergiant } from '../../../shared/supergiant/supergiant.service';
import { Notifications } from '../../../shared/notifications/notifications.service';
import { ChartsModule, BaseChartDirective } from 'ng2-charts';
import { ContextMenuService, ContextMenuComponent } from 'ngx-contextmenu';

@Component({
  selector: 'app-list-cloud-accounts',
  templateUrl: './list-cloud-accounts.component.html',
  styleUrls: ['./list-cloud-accounts.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class ListCloudAccountsComponent implements OnInit, OnDestroy {
  public subscriptions = new Subscription();
  public hasCloudAccount = false;
  rows = [];
  selected = [];
  columns = [
    { prop: 'name' },
    { prop: 'provider' },
  ];
  public displayCheck: boolean;

  private rawEvent: any;
  private contextmenuRow: any;
  private contextmenuColumn: any;
  @ViewChild(ContextMenuComponent) public basicMenu: ContextMenuComponent;

  constructor(
    private supergiant: Supergiant,
    private notifications: Notifications,
    private contextMenuService: ContextMenuService,
  ) { }

  ngOnInit() {
    this.getCloudAccounts();
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  onSelect({ selected }) {
    this.selected.splice(0, this.selected.length);
    this.selected.push(...selected);
  }

  getCloudAccounts() {
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.CloudAccounts.get()).subscribe(
      (accounts) => {
        if (Object.keys(accounts).length > 0) {
          this.hasCloudAccount = true;
          this.rows = accounts.map(account => ({
            name: account.name, provider: account.provider
          }));

          // Maintain selection of accounts:
          const selected: Array<any> = [];
          this.selected.forEach((account, index) => {
            for (const row of this.rows) {
              if (row.name === account.name) {
                selected.push(row);
                break;
              }
            }
          });
          this.selected = selected;
        }
      },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  delete() {
    if (this.selected.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No Cloud Account Selected.');
    } else {
      for (const account of this.selected) {
        this.subscriptions.add(this.supergiant.CloudAccounts.delete(account.name).subscribe(
          (data) => {
            this.notifications.display('success', 'Cloud Account: ' + account.name, 'Deleted...');
            this.selected = [];
          },
          (err) => {
            this.notifications.display('error', 'Cloud Account: ' + account.name, 'Error:' + err);
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

  contextDelete(item) {
    for (const row of this.rows) {
      if (row.id === item.id) {
        this.selected.push(row);
        this.delete();
        break;
      }
    }
  }

}

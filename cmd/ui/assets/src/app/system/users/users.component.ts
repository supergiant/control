import { Component, OnInit, OnDestroy, ViewEncapsulation, ViewChild } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { Observable } from 'rxjs/Observable';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { UsersModel } from './users.model';
import { ChartsModule, BaseChartDirective } from 'ng2-charts';
import { ContextMenuService, ContextMenuComponent } from 'ngx-contextmenu';

@Component({
  selector: 'app-users',
  templateUrl: './users.component.html',
  styleUrls: ['./users.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class Users2000Component implements OnInit, OnDestroy {
  public rows = [];
  public selected = [];
  public users = [];
  public columns = [
    { prop: 'username' },
    { prop: 'role' },
    { prop: 'token' },
  ];
  public displayCheck: boolean;

  private subscriptions = new Subscription();
  private username: string;
  private password: string;
  private role: string;
  private token: string;
  private userModel = new UsersModel;

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
    this.get();
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  onSelect({ selected }) {
    this.selected.splice(0, this.selected.length);
    this.selected.push(...selected);
  }

  get() {
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.Users.get()).subscribe(
        (users) => {
          this.rows = users.items.map(user => ({
            id: user.id, username: user.username, role: user.role, token: user.api_token
          }));

        // Maintain selection of users:
        const selected: Array<any> = [];
        this.selected.forEach((user, index) => {
          for (const row of this.rows) {
            if (row.id === user.id) {
              selected.push(row);
              break;
            }
          }
        });
        this.selected = selected;
      },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  save() {
    this.userModel.user.model.username = this.username;
    this.userModel.user.model.password = this.password;
    this.userModel.user.model.role = this.role;
    this.subscriptions.add(this.supergiant.Users.create(this.userModel.user.model).subscribe(
      (success) => {
        this.notifications.display('success', 'User: ' + this.username, 'Created...');
        this.get();
        this.username = '';
        this.password = '';
        this.role = '';
      },
      (err) => { this.notifications.display('error', 'User Create Error:', err); },
    ));
  }

  delete() {
    if (this.selected.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No User Selected.');
    } else {
      for (const user of this.selected) {
        this.subscriptions.add(this.supergiant.Users.delete(user.id).subscribe(
          (data) => {
            this.notifications.display('success', 'User: ' + user.name, 'Deleted...');
            this.selected = [];
          },
          (err) => {
            this.notifications.display('error', 'User: ' + user.name, 'Error:' + err);
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


import {timer as observableTimer,  Subscription ,  Observable } from 'rxjs';

import {switchMap} from 'rxjs/operators';
import { Component, OnInit, OnDestroy, ViewEncapsulation, ViewChild } from '@angular/core';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { ChartsModule, BaseChartDirective } from 'ng2-charts';
import { ContextMenuService, ContextMenuComponent } from 'ngx-contextmenu';

export class RepoModel {
  repo = {
    'model': {
      'name': '',
      'url': ''
    }
  };
}

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class MainComponent implements OnInit, OnDestroy {

  public rows = [];
  public selected = [];
  public columns = [
    { prop: 'name' },
    { prop: 'url' },
  ];
  public displayCheck: boolean;
  private subscriptions = new Subscription();
  private name: string;
  private url: string;
  private repoModel = new RepoModel;

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
    this.subscriptions.add(observableTimer(0, 5000).pipe(
      switchMap(() => this.supergiant.HelmRepos.get())).subscribe(
      (repos) => {
        this.rows = repos.items.map(repo => ({
          id: repo.id, name: repo.name, url: repo.url
        }));

        // Maintain selection of Helm repos:
        const selected: Array<any> = [];
        this.selected.forEach((repo, index) => {
          for (const row of this.rows) {
            if (row.id === repo.id) {
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
    this.repoModel.repo.model.name = this.name;
    this.repoModel.repo.model.url = this.url;
    this.subscriptions.add(this.supergiant.HelmRepos.create(this.repoModel.repo.model).subscribe(
      (success) => {
        this.notifications.display('success', 'Repo: ' + this.name, 'Created...');
        this.get();
        this.name = '';
        this.url = '';
      },
      (err) => { this.notifications.display('error', 'Create Error:', err); },
    ));
  }

  delete() {
    if (this.selected.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No Repo Selected.');
    } else {
      for (const repo of this.selected) {
        this.subscriptions.add(this.supergiant.HelmRepos.delete(repo.id).subscribe(
          (data) => {
            this.notifications.display('success', 'Repo: ' + repo.name, 'Deleted...');
            this.selected = [];
          },
          (err) => {
            this.notifications.display('error', 'Repo: ' + repo.name, 'Error:' + err);
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

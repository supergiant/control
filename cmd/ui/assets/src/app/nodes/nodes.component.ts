import { Component, OnDestroy, OnInit } from '@angular/core';
import { NodesService } from './nodes.service';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../shared/supergiant/supergiant.service';
import { Notifications } from '../shared/notifications/notifications.service';
import { Observable } from 'rxjs/Rx';


@Component({
  selector: 'app-nodes',
  templateUrl: './nodes.component.html',
  styleUrls: ['./nodes.component.css']
})
export class NodesComponent implements OnInit, OnDestroy {
  private nodes = [];
  subscriptions = new Subscription();

  constructor(
    private nodesService: NodesService,
    private supergiant: Supergiant,
    private notifications: Notifications,
  ) { }


  ngOnInit() {
    this.getAccounts();
  }

  getAccounts() {
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.Nodes.get()).subscribe(
      (nodes) => { this.nodes = nodes.items; },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }
}

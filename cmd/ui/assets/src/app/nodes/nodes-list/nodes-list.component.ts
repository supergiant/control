import { Component, OnDestroy, OnInit } from '@angular/core';
import { NodesService } from '../nodes.service';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { Observable } from 'rxjs/Observable';

@Component({
  selector: 'app-nodes-list',
  templateUrl: './nodes-list.component.html',
  styleUrls: ['./nodes-list.component.scss']
})
export class NodesListComponent implements OnInit, OnDestroy {
  public p: number[] = [];
  public nodes = [];
  private subscriptions = new Subscription();
  public i: number;
  public id: number;

  constructor(
    public nodesService: NodesService,
    private supergiant: Supergiant,
    private notifications: Notifications,
  ) { }


  ngOnInit() {
    this.getNodes();
  }

  getNodes() {
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.Nodes.get()).subscribe(
      (nodes) => { this.nodes = nodes.items; },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }
}

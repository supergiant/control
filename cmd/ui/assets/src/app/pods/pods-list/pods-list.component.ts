import { Component, OnDestroy, OnInit } from '@angular/core';
import { PodsService } from '../pods.service';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { Observable } from 'rxjs/Observable';


@Component({
  selector: 'app-pods-list',
  templateUrl: './pods-list.component.html',
  styleUrls: ['./pods-list.component.css']
})
export class PodsListComponent implements OnInit, OnDestroy {
  public p: number[] = [];
  public pods = [];
  private subscriptions = new Subscription();
  public i: number;
  public id: number;

  constructor(
    public podsService: PodsService,
    private supergiant: Supergiant,
    private notifications: Notifications,
  ) { }


  ngOnInit() {
    this.getPods();
  }

  getPods() {
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.KubeResources.get()).subscribe(
      (pods) => { this.pods = pods.items.filter(resource => resource.kind === 'Pod'); },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }
}

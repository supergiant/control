import { Component, OnDestroy, OnInit } from '@angular/core';
import { KubesService } from '../kubes.service';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { Observable } from 'rxjs/Observable';

@Component({
  selector: 'app-kubes-list',
  templateUrl: './kubes-list.component.html',
  styleUrls: ['./kubes-list.component.css']
})
export class KubesListComponent implements OnInit, OnDestroy {

  private kubes = [];
  subscriptions = new Subscription();

  constructor(
    private kubesService: KubesService,
    private supergiant: Supergiant,
    private notifications: Notifications,
  ) { }


  ngOnInit() {
    this.getAccounts();
  }

  getAccounts() {
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.Kubes.get()).subscribe(
      (kubes) => { this.kubes = kubes.items; },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

}

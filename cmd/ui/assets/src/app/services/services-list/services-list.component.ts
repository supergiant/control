import { Component, OnDestroy, OnInit } from '@angular/core';
import { ServicesService } from '../services.service';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { Observable } from 'rxjs/Observable';

@Component({
  selector: 'app-services-list',
  templateUrl: './services-list.component.html',
  styleUrls: ['./services-list.component.css']
})
export class ServicesListComponent implements OnInit, OnDestroy {
  p: number[] = [];
  private services = [];
  subscriptions = new Subscription();

  constructor(
    private servicesService: ServicesService,
    private supergiant: Supergiant,
    private notifications: Notifications,
  ) { }

  ngOnInit() {
    this.getServices();
  }

  getServices() {
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.KubeResources.get()).subscribe(
      (services) => { this.services = services.items.filter(resource => resource.kind === 'Service'); },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }
}


import {timer as observableTimer,  Subscription ,  Observable } from 'rxjs';

import {switchMap} from 'rxjs/operators';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { LoadBalancersService } from '../load-balancers.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';

@Component({
  selector: 'app-load-balancers-list',
  templateUrl: './load-balancers-list.component.html',
  styleUrls: ['./load-balancers-list.component.scss']
})
export class LoadBalancersListComponent implements OnInit, OnDestroy {
  public p: number[] = [];
  public loadBalancers = [];
  private subscriptions = new Subscription();
  public i: number;
  public id: number;

  constructor(
    public loadBalancersService: LoadBalancersService,
    private supergiant: Supergiant,
    private notifications: Notifications,
  ) { }

  ngOnInit() {
    this.getLoadBalancers();
  }

  getLoadBalancers() {
    this.subscriptions.add(observableTimer(0, 5000).pipe(
      switchMap(() => this.supergiant.KubeResources.get())).subscribe(
      (resources) => {
        this.loadBalancers = resources.items.filter(
          service => service.resource.spec.type === 'LoadBalancer'
        );
      },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }
}

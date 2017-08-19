import { Component, OnDestroy, OnInit } from '@angular/core';
import { LoadBalancersService } from './load-balancers.service';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../shared/supergiant/supergiant.service'
import { Notifications } from '../shared/notifications/notifications.service'
import { Observable } from 'rxjs/Rx';


@Component({
  selector: 'app-load-balancers',
  templateUrl: './load-balancers.component.html',
  styleUrls: ['./load-balancers.component.css']
})
export class LoadBalancersComponent implements OnInit {
  private loadBalancers = [];
  subscriptions = new Subscription();

  constructor(
    private loadBalancersService: LoadBalancersService,
    private supergiant: Supergiant,
    private notifications: Notifications,
  ) { }

  //get accouts when page loads
  ngOnInit() {
    this.getAccounts()
  }
  //get accounts
  getAccounts() {
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.LoadBalancers.get()).subscribe(
      (loadBalancers) => { this.loadBalancers = loadBalancers.items },
      (err) => { this.notifications.display("Warning!", "Cannot connect to Load Balancers API.", err) }))
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe()
  }
}

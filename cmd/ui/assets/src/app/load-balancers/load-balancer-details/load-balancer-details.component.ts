import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { ActivatedRoute, Router } from '@angular/router';
import { Observable } from 'rxjs/Observable';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';

@Component({
  selector: 'app-load-balancer-details',
  templateUrl: './load-balancer-details.component.html',
  styleUrls: ['./load-balancer-details.component.css']
})
export class LoadBalancerDetailsComponent implements OnInit, OnDestroy {
  id: number;
  subscriptions = new Subscription();
  loadBalancer: any;
  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private supergiant: Supergiant,
    private notifications: Notifications,
  ) { }

  ngOnInit() {
    this.id = this.route.snapshot.params.id;
    this.getLoadBalancer();
  }

  getLoadBalancer() {
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.LoadBalancers.get(this.id)).subscribe(
      (loadBalancer) => { this.loadBalancer = loadBalancer; },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  goBack() {
    this.router.navigate(['/users']);
  }
  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

}

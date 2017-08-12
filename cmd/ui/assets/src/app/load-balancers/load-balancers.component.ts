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
  private subscription: Subscription;

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
    this.subscription = Observable.timer(0, 5000)
    .switchMap(() => this.supergiant.LoadBalancers.get()).subscribe(
      (loadBalancersObj) => {
        // Because of the check boxes we must reconsile the array.
        // If does not exist locally push it locally.
        for(let loadBalancer of loadBalancersObj.json().items) {
          var present = false
          for(let uiLoadBalancer of this.loadBalancers) {
            if ( loadBalancer.id === uiLoadBalancer.id ) {present = true}
          }
          if (!present) {this.loadBalancers.push(loadBalancer)}
         }

         // If does not exist on the API remove it locally.
         for(let uiLoadBalancer of this.loadBalancers) {
           var present = false
           for(let loadBalancer of loadBalancersObj.json().items) {
             if ( loadBalancer.id === uiLoadBalancer.id ) {present = true}
           }
           if (!present) {
             var index = this.loadBalancers.indexOf(uiLoadBalancer)
             this.loadBalancers.splice(index, 1)}
          }
      },
      (err) => { this.notifications.display("warn", "Connection Issue.", err)});
  }

  ngOnDestroy(){
    this.subscription.unsubscribe();
  }
}

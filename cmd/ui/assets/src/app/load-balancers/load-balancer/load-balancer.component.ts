import { Component, Input } from '@angular/core';
import { LoadBalancersService } from '../load-balancers.service';

@Component({
  selector: '[app-load-balancer]', // tslint:disable-line
  templateUrl: './load-balancer.component.html',
  styleUrls: ['./load-balancer.component.css']
})
export class LoadBalancerComponent {
  @Input() loadBalancer: any;
  constructor(public loadBalancersService: LoadBalancersService) { }

  status(loadBalancer) {
    if (loadBalancer.status && loadBalancer.status.error && loadBalancer.status.retries === loadBalancer.status.max_retries) {
      return 'status status-danger';
    } else if (loadBalancer.status) {
      return 'status status-transitioning';
    } else if (loadBalancer.passive_status && !loadBalancer.passive_status_okay) {
      return 'status status-warning';
    } else {
      return 'status status-ok';
    }
  }
}

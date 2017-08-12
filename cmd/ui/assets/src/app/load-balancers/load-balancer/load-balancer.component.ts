import { Component, Input} from '@angular/core';
import { LoadBalancersService } from '../load-balancers.service';

@Component({
  selector: '[app-load-balancer]',
  templateUrl: './load-balancer.component.html',
  styleUrls: ['./load-balancer.component.css']
})
export class LoadBalancerComponent {
  @Input() loadBalancer: any;
  constructor(private loadBalancersService: LoadBalancersService) { }
}

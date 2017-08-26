import { Component, Input } from '@angular/core';
import { ServicesService } from '../services.service';

@Component({
  selector: '[app-service]',  // tslint:disable-line
  templateUrl: './service.component.html',
  styleUrls: ['./service.component.css']
})
export class ServiceComponent {
  @Input() service: any;
  constructor(private servicesService: ServicesService) { }

  status(service) {
    if (service.status && service.status.error && service.status.retries === service.status.max_retries) {
      return 'status status-danger';
    } else if (service.status) {
      return 'status status-transitioning';
    } else if (service.passive_status && !service.passive_status_okay) {
      return 'status status-warning';
    } else {
      return 'status status-ok';
    }
  }
}

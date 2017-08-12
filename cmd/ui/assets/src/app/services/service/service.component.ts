import { Component, Input} from '@angular/core';
import { ServicesService } from '../services.service';

@Component({
  selector: '[app-service]',
  templateUrl: './service.component.html',
  styleUrls: ['./service.component.css']
})
export class ServiceComponent {
  @Input() service: any;
  constructor(private servicesService: ServicesService) { }
}

import { Component, Input } from '@angular/core';
import { AppsService } from '../apps.service';

@Component({
  selector: '[app-deployment]',
  templateUrl: './deployment.component.html',
  styleUrls: ['./deployment.component.css']
})
export class DeploymentComponent {
  @Input() deployment: any;
  constructor(private appsService: AppsService) { }
  status(deployment) {
    if (deployment.status_value === 'DEPLOYED') {
      return 'status status-ok';
    } else {
      return 'status status-transitioning';
    }
  }
}

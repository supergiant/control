import { Component, Input } from '@angular/core';
import { KubesService } from '../kubes.service';
import { Notifications } from '../../shared/notifications/notifications.service';


@Component({
  selector: '[app-kube]',
  templateUrl: './kube.component.html',
  styleUrls: ['./kube.component.css']
})
export class KubeComponent {
  @Input() kube: any;
  constructor(
    private kubesService: KubesService,
    private notifications: Notifications,
  ) { }

  status(kube) {
    if (kube.status && kube.status.error && kube.status.retries === kube.status.max_retries) {
      return 'status status-danger';
    } else if (kube.status) {
      return 'status status-transitioning';
    } else if (kube.passive_status && !kube.passive_status_okay) {
      return 'status status-warning';
    } else {
      return 'status status-ok';
    }
  }
}

import { Component, Input} from '@angular/core';
import { KubesService } from '../kubes.service';
import { Notifications } from '../../shared/notifications/notifications.service'


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
    switch (kube.buildStatus) {
      case "completed": {
        return "status status-ok"
      }
      case "provisioning": {
        return "status status-transitioning"
      }
      case "deleting": {
        return "status status-transitioning"
      }
       default: {
         return "status status-warning"
       }
    }
  }
}

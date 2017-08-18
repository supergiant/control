import { Component, Input} from '@angular/core';
import { KubesService } from '../kubes.service';

@Component({
  selector: '[app-kube]',
  templateUrl: './kube.component.html',
  styleUrls: ['./kube.component.css']
})
export class KubeComponent {
  @Input() kube: any;
  constructor(private kubesService: KubesService) { }
  status(kube) {
    switch (kube.status) {
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

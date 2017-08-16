import { Component, Input} from '@angular/core';
import { AppsService } from '../apps.service';

@Component({
  selector: '[app-helm-app]',
  templateUrl: './helm-app.component.html',
  styleUrls: ['./helm-app.component.css']
})
export class HelmAppComponent {
  @Input() app: any;
  constructor(private appsService: AppsService) { }

}

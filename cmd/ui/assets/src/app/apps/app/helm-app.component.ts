import { Component, Input } from '@angular/core';
import { AppsService } from '../apps.service';

@Component({
  selector: '[app-helm-app]', // tslint:disable-line
  templateUrl: './helm-app.component.html',
  styleUrls: ['./helm-app.component.scss']
})
export class HelmAppComponent {
  @Input() app: any;
  constructor(public appsService: AppsService) { }

}

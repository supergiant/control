import { Component, Input } from '@angular/core';
import { SessionsService } from '../sessions.service';

@Component({
  selector: '[app-session]', // tslint:disable-line
  templateUrl: './session.component.html',
  styleUrls: ['./session.component.scss']
})
export class SessionComponent {
  @Input() session: any;
  constructor(public sessionsService: SessionsService) { }
}

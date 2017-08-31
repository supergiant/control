import { Component, Input } from '@angular/core';
import { SessionsService } from '../sessions.service';

@Component({
  selector: '[app-session]', // tslint:disable-line
  templateUrl: './session.component.html',
  styleUrls: ['./session.component.css']
})
export class SessionComponent {
  @Input() session: any;
  constructor(public sessionsService: SessionsService) { }
}

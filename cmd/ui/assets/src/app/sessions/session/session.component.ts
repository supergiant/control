import { Component, Input} from '@angular/core';
import { SessionsService } from '../sessions.service';

@Component({
  selector: '[app-session]',
  templateUrl: './session.component.html',
  styleUrls: ['./session.component.css']
})
export class SessionComponent {
  @Input() session: any;
  constructor(private sessionsService: SessionsService) { }
}

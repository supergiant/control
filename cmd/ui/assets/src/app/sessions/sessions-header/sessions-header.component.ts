import { Component } from '@angular/core';
import { SessionsService } from '../sessions.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service'
import {SessionsComponent} from '../sessions.component'
import { Subscription } from 'rxjs/Subscription';
import { Notifications } from '../../shared/notifications/notifications.service'
import { LoginComponent } from '../../login/login.component';

@Component({
  selector: 'app-sessions-header',
  templateUrl: './sessions-header.component.html',
  styleUrls: ['./sessions-header.component.css']
})
export class SessionsHeaderComponent {
  private cloudAccountsSub: Subscription;
  sessionsObj: any;

  constructor(
    private sessionsService: SessionsService,
    private sessionsComponent: SessionsComponent,
    private supergiant: Supergiant,
    private notifications: Notifications,
    private loginComponent: LoginComponent,
    ) {}

  // If the delete button is hit, the seleted sessions are deleted.
  deleteSession() {
    var selectedItems = this.sessionsService.returnSelectedSessions()
    if (selectedItems.length === 0) {
      this.notifications.display("warn", "Warning:", "No Session Selected.")
    } else {
    for(let session of selectedItems){
      this.supergiant.Sessions.delete(session.id).subscribe(
        (data) => {
            this.notifications.display("success", "Session: " + session.id, "Deleted...")
            this.sessionsComponent.getAccounts()},
        (err) => {
            this.notifications.display("error", "Session: " + session.id, "Error:" + err)},
      );
    }
  }
  }
}

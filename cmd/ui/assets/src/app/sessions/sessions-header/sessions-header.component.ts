import { Component, OnDestroy } from '@angular/core';
import { SessionsService } from '../sessions.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Subscription } from 'rxjs/Subscription';
import { Notifications } from '../../shared/notifications/notifications.service';
import { LoginComponent } from '../../login/login.component';
import { SystemModalService } from '../../shared/system-modal/system-modal.service';


@Component({
  selector: 'app-sessions-header',
  templateUrl: './sessions-header.component.html',
  styleUrls: ['./sessions-header.component.css']
})
export class SessionsHeaderComponent implements OnDestroy {
  subscriptions = new Subscription();
  sessionsObj: any;

  constructor(
    private sessionsService: SessionsService,
    private supergiant: Supergiant,
    private notifications: Notifications,
    private loginComponent: LoginComponent,
    private systemModalService: SystemModalService,
  ) { }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  openSystemModal(message) {
    this.systemModalService.openSystemModal(message);
  }

  // If the delete button is hit, the seleted sessions are deleted.
  deleteSession() {
    const selectedItems = this.sessionsService.returnSelectedSessions();
    if (selectedItems.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No Session Selected.');
    } else {
      for (const session of selectedItems) {
        this.subscriptions.add(this.subscriptions.add(this.supergiant.Sessions.delete(session.id).subscribe(
          (data) => {
            this.notifications.display('success', 'Session: ' + session.id, 'Deleted...');
          },
          (err) => {
            this.notifications.display('error', 'Session: ' + session.id, 'Error:' + err);
          },
        )));
      }
    }
  }
}

import { Component, OnInit } from '@angular/core';
import { SessionsService } from './sessions.service';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../shared/supergiant/supergiant.service'
import { Notifications } from '../shared/notifications/notifications.service'


@Component({
  selector: 'app-cloud-accounts',
  templateUrl: './sessions.component.html',
  styleUrls: ['./sessions.component.css']
})
export class SessionsComponent implements OnInit {
  sessions: any;
  private subscription: Subscription;

  constructor(
    private sessionsService: SessionsService,
    private supergiant: Supergiant,
    private notifications: Notifications,
  ) { }

  //get accouts when page loads
  ngOnInit() {
    this.getAccounts()
  }
  //get accounts
  getAccounts() {
    this.subscription = this.supergiant.Sessions.get().subscribe(
      session=>{ this.sessions = session.json().items; },
      (err) =>{ this.notifications.display("warn", "Connection Issue.", err)})
  }
}

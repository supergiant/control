import { Component, OnDestroy, OnInit } from '@angular/core';
import { SessionsService } from '../sessions.service';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { Observable } from 'rxjs/Observable';

@Component({
  selector: 'app-sessions-list',
  templateUrl: './sessions-list.component.html',
  styleUrls: ['./sessions-list.component.css']
})
export class SessionsListComponent implements OnInit, OnDestroy {
  p: number[] = [];
  sessions: any;
  private subscriptions = new Subscription();

  constructor(
    private sessionsService: SessionsService,
    private supergiant: Supergiant,
    private notifications: Notifications,
  ) { }

  ngOnInit() {
    this.getSessions();
  }

  getSessions() {
    this.subscriptions.add(Observable.timer(0, 10000)
      .switchMap(() => this.supergiant.Sessions.get()).subscribe(
      (sessions) => { this.sessions = sessions.items; },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }
  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }
}

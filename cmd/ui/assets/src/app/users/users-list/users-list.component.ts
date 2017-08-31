import { Component, OnDestroy, OnInit } from '@angular/core';
import { UsersService } from '../users.service';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { Observable } from 'rxjs/Observable';

@Component({
  selector: 'app-users-list',
  templateUrl: './users-list.component.html',
  styleUrls: ['./users-list.component.css']
})
export class UsersListComponent implements OnInit, OnDestroy {
  p: number[] = [];
  private users = [];
  private subscriptions = new Subscription();

  constructor(
    private supergiant: Supergiant,
    private notifications: Notifications,
    private usersService: UsersService,
  ) { }

  ngOnInit() {
    this.getUsers();
  }

  getUsers() {
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.Users.get()).subscribe(
      (users) => { this.users = users.items; },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }
}

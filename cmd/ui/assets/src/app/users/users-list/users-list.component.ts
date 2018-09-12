
import {timer as observableTimer,  Subscription ,  Observable } from 'rxjs';

import {switchMap} from 'rxjs/operators';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { UsersService } from '../users.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';

@Component({
  selector: 'app-users-list',
  templateUrl: './users-list.component.html',
  styleUrls: ['./users-list.component.scss']
})
export class UsersListComponent implements OnInit, OnDestroy {
  public p: number[] = [];
  public users = [];
  private subscriptions = new Subscription();
  public i: number;
  public id: number;

  constructor(
    private supergiant: Supergiant,
    private notifications: Notifications,
    public usersService: UsersService,
  ) { }

  ngOnInit() {
    this.getUsers();
  }

  getUsers() {
    this.subscriptions.add(observableTimer(0, 5000).pipe(
      switchMap(() => this.supergiant.Users.get())).subscribe(
      (users) => { this.users = users.items.filter(user => user.username !== 'support'); },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }
}

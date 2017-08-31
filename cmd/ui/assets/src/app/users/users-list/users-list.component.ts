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
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.Users.get()).subscribe(
      (users) => { this.users = users.items.filter(user => user.username !== 'support'); },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }
}

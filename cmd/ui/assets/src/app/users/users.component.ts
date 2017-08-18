import { Component, OnDestroy, OnInit } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../shared/supergiant/supergiant.service'
import { Notifications } from '../shared/notifications/notifications.service'
import { Observable } from 'rxjs/Rx';


@Component({
  selector: 'app-users',
  templateUrl: './users.component.html',
  styleUrls: ['./users.component.css']
})
export class UsersComponent implements OnInit {
  private users = [];
  private subscription: Subscription;

  constructor(
    private supergiant: Supergiant,
    private notifications: Notifications,
  ) { }

  ngOnInit() {
    this.getUsers()
  }

  getUsers() {
    this.subscription = Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.Users.get()).subscribe(
      (users) => { this.users = users.items },
      (err) => { this.notifications.display("warn", "Connection Issue.", err) });
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }
}

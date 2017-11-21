import { Component, OnInit, OnDestroy, ViewEncapsulation } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { Observable } from 'rxjs/Observable';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { UsersModel } from './users.model';

@Component({
  selector: 'app-users',
  templateUrl: './users.component.html',
  styleUrls: ['./users.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class Users2000Component implements OnInit, OnDestroy {


  private rows = [];
  private selected = [];
  private columns = [
    { prop: 'username' },
    { prop: 'role' },
  ];
  private subscriptions = new Subscription();
  private username: string;
  private password: string;
  private role: string;
  private userModel = new UsersModel;
  constructor(
    private supergiant: Supergiant,
    private notifications: Notifications,
  ) { }

  ngOnInit() {
    this.get();
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  onSelect({ selected }) {
    this.selected.splice(0, this.selected.length);
    this.selected.push(...selected);
  }

  get() {
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.Users.get()).subscribe(
      (users) => {
        this.rows = users.items.map(user => ({
          id: user.id, username: user.username, role: user.role
        }));

        // Copy over any kubes that happen to be currently selected.
        this.selected.forEach((user, index, array) => {
          for (const row of this.rows) {
            if (row.id === user.id) {
              array[index] = row;
            }
          }
        });
      },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  save() {
    this.userModel.user.model.username = this.username;
    this.userModel.user.model.password = this.password;
    this.userModel.user.model.role = this.role;
    this.subscriptions.add(this.supergiant.Users.create(this.userModel.user.model).subscribe(
      (success) => {
        this.notifications.display('success', 'User: ' + this.username, 'Created...');
        this.get();
        this.username = '';
        this.password = '';
        this.role = '';
      },
      (err) => { this.notifications.display('error', 'User Create Error:', err); },
    ));
  }

  delete() {
    if (this.selected.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No User Selected.');
    } else {
      for (const user of this.selected) {
        this.subscriptions.add(this.supergiant.Users.delete(user.id).subscribe(
          (data) => {
            this.notifications.display('success', 'User: ' + user.name, 'Deleted...');
            this.selected = [];
          },
          (err) => {
            this.notifications.display('error', 'User: ' + user.name, 'Error:' + err);
          },
        ));
      }
    }
  }



}

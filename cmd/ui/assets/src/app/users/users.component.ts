import { Component, OnDestroy, OnInit } from '@angular/core';
import { UsersService } from './users.service';
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
    private usersService: UsersService,
    private supergiant: Supergiant,
    private notifications: Notifications,
  ) { }

  //get accouts when page loads
  ngOnInit() {
    this.getAccounts()
  }
  //get accounts
  getAccounts() {
    this.subscription = Observable.timer(0, 5000)
    .switchMap(() => this.supergiant.Users.get()).subscribe(
      (usersObj) => {
        // Because of the check boxes we must reconsile the array.
        // If does not exist locally push it locally.
        for(let user of usersObj.json().items) {
          var present = false
          for(let uiUser of this.users) {
            if ( user.id === uiUser.id ) {present = true}
          }
          if (!present) {this.users.push(user)}
         }

         // If does not exist on the API remove it locally.
         for(let uiKube of this.users) {
           var present = false
           for(let kube of usersObj.json().items) {
             if ( kube.id === uiKube.id ) {present = true}
           }
           if (!present) {
             var index = this.users.indexOf(uiKube)
             this.users.splice(index, 1)}
          }
      },
      (err) => { this.notifications.display("warn", "Connection Issue.", err)});
  }

  ngOnDestroy(){
    this.subscription.unsubscribe();
  }
}

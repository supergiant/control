import { Component, OnDestroy, OnInit } from '@angular/core';
import { CloudAccountsService } from './cloud-accounts.service';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../shared/supergiant/supergiant.service'
import { Notifications } from '../shared/notifications/notifications.service'
import { Observable } from 'rxjs/Rx';


@Component({
  selector: 'app-cloud-accounts',
  templateUrl: './cloud-accounts.component.html',
  styleUrls: ['./cloud-accounts.component.css']
})
export class CloudAccountsComponent implements OnInit {
  private cloudAccounts = [];
  private subscription: Subscription;

  constructor(
    private cloudAccountsService: CloudAccountsService,
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
    .switchMap(() => this.supergiant.CloudAccounts.get()).subscribe(
      (cloudAccount) => {
        // Because of the check boxes we must reconsile the array.
        // If does not exist locally push it locally.
        for(let account of cloudAccount.json().items) {
          var present = false
          for(let uiAccount of this.cloudAccounts) {
            if ( account.id === uiAccount.id ) {present = true}
          }
          if (!present) {this.cloudAccounts.push(account)}
         }

         // If does not exist on the API remove it locally.
         for(let uiAccount of this.cloudAccounts) {
           var present = false
           for(let account of cloudAccount.json().items) {
             if ( account.id === uiAccount.id ) {present = true}
           }
           if (!present) {
             var index = this.cloudAccounts.indexOf(uiAccount)
             this.cloudAccounts.splice(index, 1)}
          }
      },
      (err) => { this.notifications.display("warn", "Connection Issue.", err)});
  }

  ngOnDestroy(){
    this.subscription.unsubscribe();
  }
}

import { Component, OnDestroy, OnInit } from '@angular/core';
import { KubesService } from './kubes.service';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../shared/supergiant/supergiant.service'
import { Notifications } from '../shared/notifications/notifications.service'
import { Observable } from 'rxjs/Rx';


@Component({
  selector: 'app-kubes',
  templateUrl: './kubes.component.html',
  styleUrls: ['./kubes.component.css']
})
export class KubesComponent implements OnInit {
  private kubes = [];
  private subscription: Subscription;

  constructor(
    private kubesService: KubesService,
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
    .switchMap(() => this.supergiant.Kubes.get()).subscribe(
      (kubesObj) => {
        // Because of the check boxes we must reconsile the array.
        // If does not exist locally push it locally.
        for(let kube of kubesObj.json().items) {
          var present = false
          for(let uiKube of this.kubes) {
            if ( kube.id === uiKube.id ) {present = true}
          }
          if (!present) {this.kubes.push(kube)}
         }

         // If does not exist on the API remove it locally.
         for(let uiKube of this.kubes) {
           var present = false
           for(let kube of kubesObj.json().items) {
             if ( kube.id === uiKube.id ) {present = true}
           }
           if (!present) {
             var index = this.kubes.indexOf(uiKube)
             this.kubes.splice(index, 1)}
          }
      },
      (err) => { this.notifications.display("warn", "Connection Issue.", err)});
  }

  ngOnDestroy(){
    this.subscription.unsubscribe();
  }
}

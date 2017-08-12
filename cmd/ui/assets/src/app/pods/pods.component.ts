import { Component, OnDestroy, OnInit } from '@angular/core';
import { PodsService } from './pods.service';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../shared/supergiant/supergiant.service'
import { Notifications } from '../shared/notifications/notifications.service'
import { Observable } from 'rxjs/Rx';


@Component({
  selector: 'app-pods',
  templateUrl: './pods.component.html',
  styleUrls: ['./pods.component.css']
})
export class PodsComponent implements OnInit {
  private pods = [];
  private subscription: Subscription;

  constructor(
    private podsService: PodsService,
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
    .switchMap(() => this.supergiant.KubeResources.get()).subscribe(
      (podsObj) => {
        // Because of the check boxes we must reconsile the array.
        // If does not exist locally push it locally.
        for(let pod of podsObj.json().items) {
          var present = false
          for(let uiPod of this.pods) {
            if ( pod.id === uiPod.id ) {present = true}
          }
          if (!present) { if (pod.kind === "Pod") {this.pods.push(pod)}}
         }

         // If does not exist on the API remove it locally.
         for(let uiPod of this.pods) {
           var present = false
           for(let pod of podsObj.json().items) {
             if ( pod.id === uiPod.id ) {present = true}
           }
           if (!present) {
             var index = this.pods.indexOf(uiPod)
             this.pods.splice(index, 1)}
          }
      },
      (err) => { this.notifications.display("warn", "Connection Issue.", err)});
  }

  ngOnDestroy(){
    this.subscription.unsubscribe();
  }
}

import { Component, OnDestroy, OnInit } from '@angular/core';
import { ServicesService } from './services.service';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../shared/supergiant/supergiant.service'
import { Notifications } from '../shared/notifications/notifications.service'
import { Observable } from 'rxjs/Rx';


@Component({
  selector: 'app-services',
  templateUrl: './services.component.html',
  styleUrls: ['./services.component.css']
})
export class ServicesComponent implements OnInit {
  private services = [];
  private subscription: Subscription;

  constructor(
    private servicesService: ServicesService,
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
      (servicesObj) => {
        // Because of the check boxes we must reconsile the array.
        // If does not exist locally push it locally.
        for(let service of servicesObj.json().items) {
          var present = false
          for(let uiService of this.services) {
            if ( service.id === uiService.id ) {present = true}
          }
          if (!present) { if (service.kind === "Service") {this.services.push(service)}}
         }

         // If does not exist on the API remove it locally.
         for(let uiService of this.services) {
           var present = false
           for(let service of servicesObj.json().items) {
             if ( service.id === uiService.id ) {present = true}
           }
           if (!present) {
             var index = this.services.indexOf(uiService)
             this.services.splice(index, 1)}
          }
      },
      (err) => { this.notifications.display("warn", "Connection Issue.", err)});
  }

  ngOnDestroy(){
    this.subscription.unsubscribe();
  }
}

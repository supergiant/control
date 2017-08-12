import { Component, OnDestroy, OnInit } from '@angular/core';
import { VolumesService } from './volumes.service';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../shared/supergiant/supergiant.service'
import { Notifications } from '../shared/notifications/notifications.service'
import { Observable } from 'rxjs/Rx';


@Component({
  selector: 'app-volumes',
  templateUrl: './volumes.component.html',
  styleUrls: ['./volumes.component.css']
})
export class VolumesComponent implements OnInit {
  private volumes = [];
  private subscription: Subscription;

  constructor(
    private volumesService: VolumesService,
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
      (volumesObj) => {
        // Because of the check boxes we must reconsile the array.
        // If does not exist locally push it locally.
        for(let volume of volumesObj.json().items) {
          var present = false
          for(let uiVolume of this.volumes) {
            if ( volume.id === uiVolume.id ) {present = true}
          }
          if (!present) { if (volume.kind === "Volume") {this.volumes.push(volume)}}
         }

         // If does not exist on the API remove it locally.
         for(let uiVolume of this.volumes) {
           var present = false
           for(let volume of volumesObj.json().items) {
             if ( volume.id === uiVolume.id ) {present = true}
           }
           if (!present) {
             var index = this.volumes.indexOf(uiVolume)
             this.volumes.splice(index, 1)}
          }
      },
      (err) => { this.notifications.display("warn", "Connection Issue.", err)});
  }

  ngOnDestroy(){
    this.subscription.unsubscribe();
  }
}

import { Component, OnDestroy, OnInit } from '@angular/core';
import { VolumesService } from './volumes.service';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../shared/supergiant/supergiant.service';
import { Notifications } from '../shared/notifications/notifications.service';
import { Observable } from 'rxjs/Observable';


@Component({
  selector: 'app-volumes',
  templateUrl: './volumes.component.html',
  styleUrls: ['./volumes.component.css']
})
export class VolumesComponent implements OnInit, OnDestroy {
  p: number[] = [];
  private volumes = [];
  subscriptions = new Subscription();

  constructor(
    private volumesService: VolumesService,
    private supergiant: Supergiant,
    private notifications: Notifications,
  ) { }

  ngOnInit() {
    this.getVolumes();
  }

  getVolumes() {
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.KubeResources.get()).subscribe(
      (volumes) => { this.volumes = volumes.items.filter(resource => resource.kind === 'Volume'); },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }
}

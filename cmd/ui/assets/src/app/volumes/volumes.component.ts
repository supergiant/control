import { Component, OnDestroy, OnInit } from '@angular/core';
import { VolumesService } from './volumes.service';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../shared/supergiant/supergiant.service';
import { Notifications } from '../shared/notifications/notifications.service';
import { Observable } from 'rxjs/Observable';


@Component({
  selector: 'app-volumes',
  templateUrl: './volumes.component.html',
  styleUrls: ['./volumes.component.scss']
})
export class VolumesComponent implements OnInit, OnDestroy {
  public p: number[] = [];
  public volumes = [];
  private subscriptions = new Subscription();
  public i: number;
  public id: number;

  constructor(
    public volumesService: VolumesService,
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

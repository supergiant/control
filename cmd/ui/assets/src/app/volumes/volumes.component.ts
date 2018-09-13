
import {timer as observableTimer,  Subscription ,  Observable } from 'rxjs';

import {switchMap} from 'rxjs/operators';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { VolumesService } from './volumes.service';
import { Supergiant } from '../shared/supergiant/supergiant.service';
import { Notifications } from '../shared/notifications/notifications.service';


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
    this.subscriptions.add(observableTimer(0, 5000).pipe(
      switchMap(() => this.supergiant.KubeResources.get())).subscribe(
      (volumes) => { this.volumes = volumes.items.filter(resource => resource.kind === 'Volume'); },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }
}

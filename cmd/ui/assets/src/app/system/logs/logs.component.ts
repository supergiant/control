
import {timer as observableTimer,  Subscription ,  Observable } from 'rxjs';

import {switchMap} from 'rxjs/operators';
import { Component, OnInit, OnDestroy, AfterViewInit } from '@angular/core';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { SystemModalService } from '../../shared/system-modal/system-modal.service';

@Component({
  selector: 'app-logs',
  templateUrl: './logs.component.html',
  styleUrls: ['./logs.component.scss']
})
export class LogsComponent implements OnInit, OnDestroy, AfterViewInit {
  subscriptions = new Subscription();
  private logData: any;
  private notificationItem: any;
  private notificationItems = [];

  constructor(
    private supergiant: Supergiant,
    private notifications: Notifications,
    private systemModalService: SystemModalService,
  ) { }

  ngAfterViewInit() {
    this.notificationItems = this.systemModalService.notifications;
  }

  ngOnInit() {
    this.subscriptions.add(observableTimer(0, 1000).pipe(
      switchMap(() => this.supergiant.Logs.get())).subscribe(
        (data) => {
          this.logData = data;
          this.logData = this.logData.replace(/[\x00-\x7F]\[\d+mINFO[\x00-\x7F]\[0m/g, '<span class=\'text-info\'>INFO</span> ');
          this.logData = this.logData.replace(/[\x00-\x7F]\[\d+mWARN[\x00-\x7F]\[0m/g, '<span class=\'text-warning\'>WARN</span> ');
          this.logData = this.logData.replace(/[\x00-\x7F]\[\d+mERRO[\x00-\x7F]\[0m/g, '<span class=\'text-danger\'>ERRO</span> ');
          this.logData = this.logData.replace(/[\x00-\x7F]\[\d+mDEBU[\x00-\x7F]\[0m/g, '<span class=\'text-muted\'>DEBU</span> ');
        },
    ));
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }
}

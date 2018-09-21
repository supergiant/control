
import {timer as observableTimer,  Subscription ,  Observable } from 'rxjs';

import {switchMap} from 'rxjs/operators';
import { Component, OnInit, OnDestroy, ViewChild, ViewEncapsulation } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { ChartsModule, BaseChartDirective } from 'ng2-charts';
import { SystemModalService } from '../../shared/system-modal/system-modal.service';
import { convertIsoToHumanReadable } from '../../shared/helpers/helpers';
import { LoginComponent } from '../../login/login.component';
import { Location } from '@angular/common';

@Component({
  selector: 'app-pod-details',
  templateUrl: './pod-details.component.html',
  styleUrls: ['./pod-details.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class PodDetailsComponent implements OnInit, OnDestroy {
  id: number;
  subscriptions = new Subscription();
  pod: any;
  podString: string;
  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private supergiant: Supergiant,
    private notifications: Notifications,
    private chartsModule: ChartsModule,
    private systemModalService: SystemModalService,
    public loginComponent: LoginComponent,
    private location: Location,
  ) { }

  // CPU Usage
  public cpuChartData: Array<any> = [];
  public cpuChartOptions: any = {
    responsive: true
  };
  public cpuChartLabels: Array<any> = [];
  public cpuChartType = 'line';

  // RAM Usage
  public ramChartData: Array<any> = [];
  public ramChartOptions: any = {
    responsive: true
  };
  public ramChartLabels: Array<any> = [];
  public ramChartType = 'line';


  isDataAvailable = false;
  ngOnInit() {
    this.id = this.route.snapshot.params.id;
    this.getPod();
  }

  openSystemModal(message) {
    this.systemModalService.openSystemModal(message);
  }

  getPod() {
    this.subscriptions.add(observableTimer(0, 5000).pipe(
      switchMap(() => this.supergiant.KubeResources.get(this.id))).subscribe(
        (pod) => {
          this.pod = pod;
          if (this.pod.extra_data && this.pod.extra_data.metrics.cpu_usage) {
            this.isDataAvailable = true;
            this.cpuChartData = [
              { label: 'CPU Usage', data: this.pod.extra_data.metrics.cpu_usage.map((data) => data.value) },
              // this should be set to the length of largest array.
            ];
            this.ramChartLabels = this.pod.extra_data.metrics.cpu_usage.map((data) => convertIsoToHumanReadable(data.timestamp));

            this.ramChartData = [
              { label: 'RAM Usage', data: this.pod.extra_data.metrics.ram_usage.map((data) => data.value / 1073741824) },
              // this should be set to the length of largest array.
            ];
            this.cpuChartLabels = this.pod.extra_data.metrics.ram_usage.map((data) => convertIsoToHumanReadable(data.timestamp));
          }
          this.podString = JSON.stringify(this.pod, null, 2);
        },
        (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  padArrayWithDefault(arr: any, n: number) {
    let tmpArr = [];
    tmpArr = arr.slice(0);
    while (tmpArr.length < n) {
      let count = 0;
      arr = tmpArr.slice(0);
      arr.reduce((previous, current, index) => {
        if (previous && tmpArr.length < n) {
          const average = (current + previous) / 2;
          tmpArr.splice(index + count, 0, average);
          count += 1;
        }
        return current;
      });
    }
    return tmpArr;
  }

  goBack() {
    this.location.back();
  }
  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

}

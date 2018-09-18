
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

import "brace/mode/json";

@Component({
  selector: 'app-node-details',
  templateUrl: './node-details.component.html',
  styleUrls: ['./node-details.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class NodeDetailsComponent implements OnInit, OnDestroy {
  id: number;
  subscriptions = new Subscription();
  node: any;
  nodeString: string;
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
    this.getNode();
  }

  openSystemModal(message) {
    this.systemModalService.openSystemModal(message);
  }

  getNode() {
    this.subscriptions.add(observableTimer(0, 5000).pipe(
      switchMap(() => this.supergiant.Nodes.get(this.id))).subscribe(
        (node) => {
          this.node = node;
          if (this.node.extra_data && this.node.extra_data.cpu_usage_rate && this.node.extra_data.cpu_node_capacity) {
            this.isDataAvailable = true;
            this.cpuChartData = [
              { label: 'CPU Usage', data: this.node.extra_data.cpu_usage_rate.map((data) => data.value) },
              { label: 'CPU Capacity', data: this.node.extra_data.cpu_node_capacity.map((data) => data.value) },
              // this should be set to the length of largest array.
            ];
            this.ramChartLabels = this.node.extra_data.cpu_usage_rate.map((data) => convertIsoToHumanReadable(data.timestamp));

            this.ramChartData = [
              { label: 'RAM Usage', data: this.node.extra_data.memory_usage.map((data) => data.value / 1073741824) },
              {
                label: 'RAM Capacity',
                data: this.node.extra_data.memory_node_capacity.map((data) => data.value / 1073741824)
              },
              // this should be set to the length of largest array.
            ];
            this.cpuChartLabels = this.node.extra_data.memory_usage.map((data) => convertIsoToHumanReadable(data.timestamp));
          }
          this.nodeString = JSON.stringify(this.node, null, 2);
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

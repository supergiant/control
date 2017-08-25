import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { ActivatedRoute, Router } from '@angular/router';
import { Observable } from 'rxjs/Observable';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { ChartsModule, BaseChartDirective } from 'ng2-charts';

@Component({
  selector: 'app-kube-details',
  templateUrl: './kube-details.component.html',
  styleUrls: ['./kube-details.component.css']
})
export class KubeDetailsComponent implements OnInit, OnDestroy {
  id: number;
  subscriptions = new Subscription();
  kube: any;
  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private supergiant: Supergiant,
    private notifications: Notifications,
    private chartsModule: ChartsModule,
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
    this.subscriptions.add(this.route.params.subscribe(params => {
      this.id = +params['id']; // (+) converts string 'id' to a number
    }));
    this.getAccounts();
  }

  getAccounts() {
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.Kubes.get(this.id)).subscribe(
      (kube) => {
        this.kube = kube.items[0];
        const cpuCapacity = this.kube.extra_data.kube_cpu_capacity.value;
        this.cpuChartData = [
          { label: 'CPU Usage', data: this.kube.extra_data.cpu_usage_rate.map((data) => data.value) },
          { label: 'CPU Capacity', data: this.padArrayWithDefault([cpuCapacity, cpuCapacity], 45) },
          // this should be set to the length of largest array.
        ];
        const ramCapacity = this.kube.extra_data.kube_memory_capacity.value / 1073741824;
        this.ramChartData = [
          { label: 'RAM Usage', data: this.kube.extra_data.memory_usage.map((data) => data.value / 1073741824) },
          {
            label: 'CPU Capacity',
            data: this.padArrayWithDefault([ramCapacity, ramCapacity], 45)
          },
          // this should be set to the length of largest array.
        ];
        this.cpuChartLabels = this.kube.extra_data.memory_usage.map((data) => data.timestamp);
        this.ramChartLabels = this.kube.extra_data.cpu_usage_rate.map((data) => data.timestamp);
        this.isDataAvailable = true;
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
    this.router.navigate(['/kubes']);
  }
  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

}

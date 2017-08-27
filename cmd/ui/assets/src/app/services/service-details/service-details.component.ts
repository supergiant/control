import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { ActivatedRoute, Router } from '@angular/router';
import { Observable } from 'rxjs/Observable';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { ChartsModule, BaseChartDirective } from 'ng2-charts';

@Component({
  selector: 'app-service-details',
  templateUrl: './service-details.component.html',
  styleUrls: ['./service-details.component.css']
})
export class ServiceDetailsComponent implements OnInit, OnDestroy {
  id: number;
  subscriptions = new Subscription();
  service: any;
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
    this.id = this.route.snapshot.params.id;
    this.getService();
  }

  getService() {
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.KubeResources.get(this.id)).subscribe(
      (service) => {
        this.service = service;
        if (this.service.extra_data) {
          this.isDataAvailable = true;
          this.cpuChartData = [
            { label: 'CPU Usage', data: this.service.extra_data.metrics.cpu_usage.map((data) => data.value) },
            // this should be set to the length of largest array.
          ];
          this.ramChartLabels = this.service.extra_data.metrics.cpu_usage.map((data) => data.timestamp);

          this.ramChartData = [
            { label: 'RAM Usage', data: this.service.extra_data.metrics.ram_usage.map((data) => data.value / 1073741824) },
            // this should be set to the length of largest array.
          ];
          this.cpuChartLabels = this.service.extra_data.metrics.ram_usage.map((data) => data.timestamp);
        }

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
    this.router.navigate(['/services']);
  }
  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

}

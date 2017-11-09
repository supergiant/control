import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { ActivatedRoute, Router } from '@angular/router';
import { Observable } from 'rxjs/Observable';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { ChartsModule, BaseChartDirective } from 'ng2-charts';
import { SystemModalService } from '../../shared/system-modal/system-modal.service';
import { LoginComponent } from '../../login/login.component';

@Component({
  selector: 'app-node-details',
  templateUrl: './node-details.component.html',
  styleUrls: ['./node-details.component.scss']
})
export class NodeDetailsComponent implements OnInit, OnDestroy {
  id: number;
  subscriptions = new Subscription();
  node: any;
  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private supergiant: Supergiant,
    private notifications: Notifications,
    private chartsModule: ChartsModule,
    private systemModalService: SystemModalService,
    public loginComponent: LoginComponent,
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
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.Nodes.get(this.id)).subscribe(
      (node) => {
        this.node = node;
        if (this.node.extra_data && this.node.extra_data.cpu_usage_rate && this.node.extra_data.cpu_node_capacity) {
          this.isDataAvailable = true;
          this.cpuChartData = [
            { label: 'CPU Usage', data: this.node.extra_data.cpu_usage_rate.map((data) => data.value) },
            { label: 'CPU Capacity', data: this.node.extra_data.cpu_node_capacity.map((data) => data.value) },
            // this should be set to the length of largest array.
          ];
          this.ramChartLabels = this.node.extra_data.cpu_usage_rate.map((data) => data.timestamp);

          this.ramChartData = [
            { label: 'RAM Usage', data: this.node.extra_data.memory_usage.map((data) => data.value / 1073741824) },
            {
              label: 'RAM Capacity',
              data: this.node.extra_data.memory_node_capacity.map((data) => data.value / 1073741824)
            },
            // this should be set to the length of largest array.
          ];
          this.cpuChartLabels = this.node.extra_data.memory_usage.map((data) => data.timestamp);
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
    this.router.navigate(['/nodes']);
  }
  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

}

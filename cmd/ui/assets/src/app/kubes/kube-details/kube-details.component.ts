import { Component, OnInit, OnDestroy } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { ActivatedRoute, Router } from '@angular/router';
import { Observable } from 'rxjs/Observable';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { ChartsModule } from 'ng2-charts';

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
  cpuUsage: number[];
  cpuUsageTimestamps: string[];
  public lineChartData: Array<any> = [{
    label: 'CPU Usage',
    data: this.cpuUsage,
  }];
  public lineChartLabels: Array<any> = [];
  public lineChartType: string = 'line';


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
        this.lineChartData = [{ label: 'CPU Usage', data: this.kube.extra_data.cpu_usage_rate.map((data) => data.value), }];
        this.lineChartLabels = this.kube.extra_data.cpu_usage_rate.map((data) => data.timestamp);
      },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  goBack() {
    this.router.navigate(['/kubes']);
  }
  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

}

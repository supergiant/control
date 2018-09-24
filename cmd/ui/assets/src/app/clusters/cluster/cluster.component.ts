import { Component, OnInit, OnDestroy, ViewChild, ElementRef, ViewContainerRef, ViewEncapsulation } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { Location } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { Observable } from 'rxjs/Observable';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { convertIsoToHumanReadable } from '../../shared/helpers/helpers';
import { DomSanitizer, SafeResourceUrl } from '@angular/platform-browser';
import { ChartsModule, BaseChartDirective } from 'ng2-charts';
import { MatTableDataSource, MatSort, MatPaginator } from '@angular/material';


@Component({
  selector: 'app-cluster',
  templateUrl: './cluster.component.html',
  styleUrls: ['./cluster.component.scss'],
  // encapsulation: ViewEncapsulation.None
})
export class ClusterComponent implements OnInit, OnDestroy {
  public approws = [];
  public lbrows = [];
  public showRaw = false;
  public hasApps = false;
  public hasLB = false;
  id: number;
  subscriptions = new Subscription();
  public kube: any;
  public kubeString: string;
  url: string;
  public isLoading: Boolean;
  public secureSrc: SafeResourceUrl;
  public planets = [];
  public planetName: string;
  public rowChartOptions: any = {
    responsive: true,
    maintainAspectRatio: true,
  };
  public rowChartColors: Array<any> = [
    { // grey
      backgroundColor: 'rgba(148,159,177,0.2)',
      borderColor: 'rgba(148,159,177,1)',
      pointBackgroundColor: 'rgba(148,159,177,1)',
      pointBorderColor: '#fff',
      pointHoverBackgroundColor: '#fff',
      pointHoverBorderColor: 'rgba(148,159,177,0.8)'
    },
    { // dark grey
      backgroundColor: 'rgba(77,83,96,0.2)',
      borderColor: 'rgba(77,83,96,1)',
      pointBackgroundColor: 'rgba(77,83,96,1)',
      pointBorderColor: '#fff',
      pointHoverBackgroundColor: '#fff',
      pointHoverBorderColor: 'rgba(77,83,96,1)'
    },
    { // grey
      backgroundColor: 'rgba(148,159,177,0.2)',
      borderColor: 'rgba(148,159,177,1)',
      pointBackgroundColor: 'rgba(148,159,177,1)',
      pointBorderColor: '#fff',
      pointHoverBackgroundColor: '#fff',
      pointHoverBorderColor: 'rgba(148,159,177,0.8)'
    }
  ];

  @ViewChild(MatSort) sort: MatSort;
  @ViewChild(MatPaginator) paginator: MatPaginator;


  // linter is angry about the boolean typing but without it charts

  constructor(
    private route: ActivatedRoute,
    private location: Location,
    private router: Router,
    private supergiant: Supergiant,
    private notifications: Notifications,
    private sanitizer: DomSanitizer,
  ) { }

  public cpuChartColors: Array<any> = [
    { // grey
      backgroundColor: 'rgba(148,159,177,0.2)',
      borderColor: 'rgba(148,159,177,1)',
      pointBackgroundColor: 'rgba(148,159,177,1)',
      pointBorderColor: '#fff',
      pointHoverBackgroundColor: '#fff',
      pointHoverBorderColor: 'rgba(148,159,177,0.8)'
    },
    { // dark grey
      backgroundColor: 'rgba(77,83,96,0.2)',
      borderColor: 'rgba(77,83,96,1)',
      pointBackgroundColor: 'rgba(77,83,96,1)',
      pointBorderColor: '#fff',
      pointHoverBackgroundColor: '#fff',
      pointHoverBorderColor: 'rgba(77,83,96,1)'
    },
    { // grey
      backgroundColor: 'rgba(148,159,177,0.2)',
      borderColor: 'rgba(148,159,177,1)',
      pointBackgroundColor: 'rgba(148,159,177,1)',
      pointBorderColor: '#fff',
      pointHoverBackgroundColor: '#fff',
      pointHoverBorderColor: 'rgba(148,159,177,0.8)'
    }
  ];

  public ramChartColors: Array<any> = [
    { // grey
      backgroundColor: 'rgba(148,159,177,0.2)',
      borderColor: 'rgba(148,159,177,1)',
      pointBackgroundColor: 'rgba(148,159,177,1)',
      pointBorderColor: '#fff',
      pointHoverBackgroundColor: '#fff',
      pointHoverBorderColor: 'rgba(148,159,177,0.8)'
    },
    { // dark grey
      backgroundColor: 'rgba(77,83,96,0.2)',
      borderColor: 'rgba(77,83,96,1)',
      pointBackgroundColor: 'rgba(77,83,96,1)',
      pointBorderColor: '#fff',
      pointHoverBackgroundColor: '#fff',
      pointHoverBorderColor: 'rgba(77,83,96,1)'
    },
    { // grey
      backgroundColor: 'rgba(148,159,177,0.2)',
      borderColor: 'rgba(148,159,177,1)',
      pointBackgroundColor: 'rgba(148,159,177,1)',
      pointBorderColor: '#fff',
      pointHoverBackgroundColor: '#fff',
      pointHoverBorderColor: 'rgba(148,159,177,0.8)'
    }
  ];

  // CPU Usage
  // I can't get this to update with push, so it has to have the same struct as eventual data.
  public cpuChartData: Array<any> = [{ data: [] }, { data: [] }];
  public cpuChartOptions: any = {
    responsive: true,
    maintainAspectRatio: true,
    scales: {
      yAxes: [{
        ticks: {
          beginAtZero: true,
        }
      }]
    }
  };
  public cpuChartLabels: Array<any> = [];
  public cpuChartType = 'line';
  public cpuChartLegend = true;

  // RAM Usage
  public ramChartData: Array<any> = [{ data: [] }, { data: [] }];
  public ramChartOptions: any = {
    responsive: true,
    maintainAspectRatio: true,
    scales: {
      yAxes: [{
        ticks: {
          beginAtZero: true,
        }
      }]
    }
  };
  public ramChartLabels: Array<any> = [];
  public ramChartType = 'line';

  isDataAvailable = false;

  // machine list vars
  machines: any;
  machineListColumns = ["role", "name", "id", "region", "public_ip"];

  ngOnInit() {
    this.id = this.route.snapshot.params.id;
    this.getKube();
  }

  usageOrZeroCPU(usage) {
    if (usage == null) {
      return ([0, 0, 0, 0, 0, 0, 0, 0, 0, 0]);
    } else {
      return usage.cpu_usage_rate.map((data) => data.value);
    }
  }


  onAppActivate(activated) {
    console.log(activated);
    if (activated.type === 'click' && activated.column.name !== 'checkbox') {
      this.router.navigate(['/clusters', activated.row.id]);
    }
  }

  // getKube() {
  //   this.subscriptions.add(Observable.timer(0, 20000)
  //     .switchMap(() => this.supergiant.Kubes.get(this.id)).subscribe(
  //       (kube) => {
  //         this.kube = kube;
  //         this.kubeString = JSON.stringify(this.kube, null, 2);
  //         if (this.kube.extra_data &&
  //           this.kube.extra_data.cpu_usage_rate &&
  //           this.kube.extra_data.kube_cpu_capacity) {
  //           this.isDataAvailable = true;
  //           this.cpuChartLabels.length = 0;
  //           let tempArray = this.kube.extra_data.cpu_usage_rate.map(
  //             (data) => convertIsoToHumanReadable(data.timestamp));
  //           for (const row of tempArray) {
  //             this.cpuChartLabels.push(row);
  //           }
  //           this.cpuChartData = [
  //             {
  //               label: 'CPU Usage',
  //               data: this.kube.extra_data.cpu_usage_rate.map((data) => data.value)
  //             },
  //             {
  //               label: 'CPU Capacity',
  //               data: this.kube.extra_data.kube_cpu_capacity.map((data) => data.value)
  //             },
  //             // this should be set to the length of largest array.
  //           ];
  //           this.ramChartLabels.length = 0;
  //           tempArray = this.kube.extra_data.memory_usage.map(
  //             (data) => convertIsoToHumanReadable(data.timestamp));
  //           for (const row of tempArray) {
  //             this.ramChartLabels.push(row);
  //           }
  //           this.ramChartData = [
  //             {
  //               label: 'RAM Usage',
  //               data: this.kube.extra_data.memory_usage.map((data) => data.value / 1073741824)
  //             },
  //             {
  //               label: 'RAM Capacity',
  //               data: this.kube.extra_data.kube_memory_capacity.map((data) => data.value / 1073741824)
  //             },
  //             // this should be set to the length of largest array.
  //           ];
  //         }

  //         this.hasApps = false;
  //         if (kube.helm_releases) {
  //           this.hasApps = true;
  //           this.approws = kube.helm_releases.map(app => ({
  //             id: app.id,
  //             name: app.name,
  //             version: app.revision,
  //             appname: app.chart_name,
  //             appversion: app.chart_version,
  //             statusvalue: app.status_value,
  //           }));
  //         }

  //         // // FAKEDATA
  //         // this.hasApps = true;
  //         // this.approws.push({id: '12345', name: 'fake-app', version: '1.2.3', appname: 'fake-wordpress',
  //         //   appversion: '3.4.5', statusvalue: 'A OK'});

  //         this.hasLB = false;
  //         if (kube.load_balancers) {
  //           this.hasLB = true;
  //           this.lbrows = kube.load_balancers.map(lb => ({
  //             id: lb.id,
  //             name: lb.name,
  //             ip: lb.ip,
  //           }));
  //         }

  //         // // FAKEDATA
  //         // this.hasLB = true;
  //         // this.lbrows.push({id: '12345', name: 'fake-lb', ip: '1.2.3.4'});
  //       },
  //       (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));

  //   // Get any planets
  //   this.subscriptions.add(Observable.timer(0, 20000)
  //     .switchMap(() => this.supergiant.KubeResources.get()).subscribe(
  //       (services) => {
  //         this.planets = services.items.filter(
  //           planet => {
  //             if (planet.resource.metadata.labels) {
  //               return planet.resource.metadata.labels['kubernetes.io/cluster-service'] === 'true';
  //             }
  //           }
  //         );
  //       },
  //       (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  // }

  getKube() {
    this.subscriptions.add(Observable.timer(0, 20000)
      .switchMap(() => this.supergiant.Kubes.get(this.id)).subscribe(
        kube => {
          this.kube = kube;
          this.machines = new MatTableDataSource(kube.masters.concat(kube.nodes));
          this.machines.sort = this.sort;
          this.machines.paginator = this.paginator;
          console.log(kube);},
        err => console.log(err)
      ))
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

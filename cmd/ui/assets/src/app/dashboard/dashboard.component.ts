import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import { Subscription } from 'rxjs';
import { Supergiant } from '../shared/supergiant/supergiant.service';
import { convertIsoToHumanReadable } from '../shared/helpers/helpers';
import { MatTableDataSource } from '@angular/material';

// TEMPORARY
import { AuthService } from '../shared/supergiant/auth/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class DashboardComponent implements OnInit {
  public subscriptions = new Subscription();
  public cloudAccounts: any;
  public hasCluster = false;
  public hasApp = false;
  public appCount = 0;
  public events: Array<any> = ['No Cluster Events (disabled in beta currently)'];
  public newses: Array<any> = ['No Recent News (disabled in beta currently)'];
  // lineChart
  public lineChartData: Array<any> = [{ data: [] }, { data: [] }];
  public lineChartLabels: Array<any> = [];
  public lineChartOptions: any = {
    responsive: true,
    scales: {
      xAxes: [{
        scaleLabel: {
          display: false
        }
      }]
    }
  };
  public lineChartColors: Array<any> = [
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
  public lineChartLegend: boolean = true;
  public lineChartType: string = 'line';

  public clusters: any;

  public clusterListView: string;

  constructor(
    private supergiant: Supergiant,
    public auth: AuthService,
    private router: Router
  ) { }

  setClusterListView(view) {
    this.clusterListView = view;
  }

  getKube(id) {
    this.subscriptions.add(this.supergiant.Kubes.get(id).subscribe(
      (kube) => {
        if (kube.extra_data &&
          kube.extra_data.cpu_usage_rate &&
          kube.extra_data.kube_cpu_capacity) {
          this.lineChartLabels.length = 0;
          let tempArray = kube.extra_data.cpu_usage_rate.map(
            (data) => convertIsoToHumanReadable(data.timestamp));
          for (const row of tempArray) {
            this.lineChartLabels.push(row);
          }

          tempArray = [
            {
              label: 'CPU Usage',
              data: kube.extra_data.cpu_usage_rate.map((data) => data.value)
            },
            {
              label: 'CPU Capacity',
              data: kube.extra_data.kube_cpu_capacity.map((data) => data.value)
            },
            // this should be set to the length of largest array.
          ];
          //linter is angry but it works, can change it later
          this.lineChartData[0]['label'] = 'CPU Usage';
          for (const i in tempArray[0]['data']) {
            const previous = Number(this.lineChartData[0]['data'][i]) || 0;
            tempArray[0]['data'][i] = previous + tempArray[0]['data'][i];
          }

          for (const i in tempArray[1]['data']) {
            const previous = Number(this.lineChartData[1]['data'][i]) || 0;
            tempArray[1]['data'][i] = previous + tempArray[1]['data'][i];
          }
          this.lineChartData = [
            { label: 'CPU Usage', data: tempArray[0]['data'] },
            { label: 'CPU Capacity', data: tempArray[1]['data'] }
          ];
        }
      }
    ))
  }

  getCloudAccounts() {
    this.subscriptions.add(this.supergiant.CloudAccounts.get().subscribe(
      (cloudAccounts) => {
        this.cloudAccounts = cloudAccounts;
      }));
  }

  getClusters() {
    this.subscriptions.add(this.supergiant.Kubes.get().subscribe(
      clusters => {
        // TODO: this is terrible, fix after demo
        clusters.map(c => c.dataSource = new MatTableDataSource([
          {
            state: c.state,
            region: c.region,
            accountName: c.accountName,
            K8SVersion: c.K8SVersion,
            masters: Object.keys(c.masters),
            nodes: Object.keys(c.nodes),
            operatingSystem: c.operatingSystem,
            dockerVersion: c.dockerVersion,
            helmVersion: c.helmVersion,
            rbacEnabled: c.rbacEnabled,
            arch: c.arch
          }]));
        this.clusters = clusters;
        if (clusters.length > 5) {
          this.clusterListView = "list"
        } else {
          this.clusterListView = "orb"
        }
     },
    err => console.error(err)
   ));
  }

  getDeployments() {
    this.subscriptions.add(this.supergiant.HelmReleases.get().subscribe(
      (deployments) => {
        if (Object.keys(deployments.items).length > 0) {
          this.hasApp = true;
          this.appCount = Object.keys(deployments.items).length;
        }
      }));
    // this.hasApp = true;
  }

  logout() {
    this.auth.logout();
    this.router.navigate(['']);
  }

  ngOnInit() {
    this.getCloudAccounts();
    this.getClusters();
    // this.getDeployments();
  }

}

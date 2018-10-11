import { of, Subscription, timer as observableTimer } from 'rxjs';

import { catchError, filter, switchMap } from 'rxjs/operators';
import { ChangeDetectionStrategy, Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { Location } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { SafeResourceUrl } from '@angular/platform-browser';
import { MatDialog, MatPaginator, MatSort, MatTableDataSource } from '@angular/material';
import { ConfirmModalComponent } from '../../shared/modals/confirm-modal/confirm-modal.component';
import { HttpClient } from '@angular/common/http';


@Component({
  selector: 'app-cluster',
  templateUrl: './cluster.component.html',
  styleUrls: [ './cluster.component.scss' ],
  changeDetection: ChangeDetectionStrategy.OnPush
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
      pointHoverBorderColor: 'rgba(148,159,177,0.8)',
    },
    { // dark grey
      backgroundColor: 'rgba(77,83,96,0.2)',
      borderColor: 'rgba(77,83,96,1)',
      pointBackgroundColor: 'rgba(77,83,96,1)',
      pointBorderColor: '#fff',
      pointHoverBackgroundColor: '#fff',
      pointHoverBorderColor: 'rgba(77,83,96,1)',
    },
    { // grey
      backgroundColor: 'rgba(148,159,177,0.2)',
      borderColor: 'rgba(148,159,177,1)',
      pointBackgroundColor: 'rgba(148,159,177,1)',
      pointBorderColor: '#fff',
      pointHoverBackgroundColor: '#fff',
      pointHoverBorderColor: 'rgba(148,159,177,0.8)',
    },
  ];

  @ViewChild(MatSort) sort: MatSort;
  @ViewChild(MatPaginator) paginator: MatPaginator;


  // linter is angry about the boolean typing but without it charts
  public cpuChartColors: Array<any> = [
    { // grey
      backgroundColor: 'rgba(148,159,177,0.2)',
      borderColor: 'rgba(148,159,177,1)',
      pointBackgroundColor: 'rgba(148,159,177,1)',
      pointBorderColor: '#fff',
      pointHoverBackgroundColor: '#fff',
      pointHoverBorderColor: 'rgba(148,159,177,0.8)',
    },
    { // dark grey
      backgroundColor: 'rgba(77,83,96,0.2)',
      borderColor: 'rgba(77,83,96,1)',
      pointBackgroundColor: 'rgba(77,83,96,1)',
      pointBorderColor: '#fff',
      pointHoverBackgroundColor: '#fff',
      pointHoverBorderColor: 'rgba(77,83,96,1)',
    },
    { // grey
      backgroundColor: 'rgba(148,159,177,0.2)',
      borderColor: 'rgba(148,159,177,1)',
      pointBackgroundColor: 'rgba(148,159,177,1)',
      pointBorderColor: '#fff',
      pointHoverBackgroundColor: '#fff',
      pointHoverBorderColor: 'rgba(148,159,177,0.8)',
    },
  ];
  public ramChartColors: Array<any> = [
    { // grey
      backgroundColor: 'rgba(148,159,177,0.2)',
      borderColor: 'rgba(148,159,177,1)',
      pointBackgroundColor: 'rgba(148,159,177,1)',
      pointBorderColor: '#fff',
      pointHoverBackgroundColor: '#fff',
      pointHoverBorderColor: 'rgba(148,159,177,0.8)',
    },
    { // dark grey
      backgroundColor: 'rgba(77,83,96,0.2)',
      borderColor: 'rgba(77,83,96,1)',
      pointBackgroundColor: 'rgba(77,83,96,1)',
      pointBorderColor: '#fff',
      pointHoverBackgroundColor: '#fff',
      pointHoverBorderColor: 'rgba(77,83,96,1)',
    },
    { // grey
      backgroundColor: 'rgba(148,159,177,0.2)',
      borderColor: 'rgba(148,159,177,1)',
      pointBackgroundColor: 'rgba(148,159,177,1)',
      pointBorderColor: '#fff',
      pointHoverBackgroundColor: '#fff',
      pointHoverBorderColor: 'rgba(148,159,177,0.8)',
    },
  ];
  // I can't get this to update with push, so it has to have the same struct as eventual data.
  public cpuChartData: Array<any> = [ { data: [] }, { data: [] } ];

  // CPU Usage
  public cpuChartOptions: any = {
    responsive: true,
    maintainAspectRatio: true,
    scales: {
      yAxes: [ {
        ticks: {
          beginAtZero: true,
        },
      } ],
    },
  };
  public cpuChartLabels: Array<any> = [];
  public cpuChartType = 'line';
  public cpuChartLegend = true;
  // RAM Usage
  public ramChartData: Array<any> = [ { data: [] }, { data: [] } ];
  public ramChartOptions: any = {
    responsive: true,
    maintainAspectRatio: true,
    scales: {
      yAxes: [ {
        ticks: {
          beginAtZero: true,
        },
      } ],
    },
  };
  public ramChartLabels: Array<any> = [];
  public ramChartType = 'line';
  isDataAvailable = false;
  // machine list vars
  machines: any;
  machineListColumns = [ 'role', 'name', 'id', 'region', 'publicIp' ];

  constructor(
    private route: ActivatedRoute,
    private location: Location,
    private router: Router,
    private supergiant: Supergiant,
    public dialog: MatDialog,
    public http: HttpClient,
  ) {
  }

  ngOnInit() {
    this.id = this.route.snapshot.params.id;
    this.getKube();
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  getKube() {
    this.subscriptions.add(observableTimer(0, 20000).pipe(
      switchMap(() => this.supergiant.Kubes.get(this.id))).subscribe(
      this.renderKube.bind(this),
      err => console.error(err),
    ));
  }


  renderKube(kube) {

    this.kube = kube;

    let masters = kube.masters;
    let nodes = kube.nodes;

    const mastersArr = Object
      .keys(masters)
      .map(name => Object.assign({}, { name }, masters[ name ]));

    const nodesArr = Object
      .keys(nodes)
      .map(name => Object.assign({}, { name }, nodes[ name ]));


    this.machines = new MatTableDataSource(mastersArr.concat(nodesArr));
    this.machines.sort = this.sort;
    this.machines.paginator = this.paginator;
  }

  goBack() {
    this.location.back();
  }

  removeNode(nodeName: string, target) {
    const dialogRef = this.initDialog(target);

    dialogRef
      .afterClosed()
      .pipe(
        filter(isConfirmed => isConfirmed),
        switchMap(() => this.supergiant.Nodes.delete(this.id, nodeName)),
        switchMap(() => this.supergiant.Kubes.get(this.id)),
        catchError((error) => of(error)),
      ).subscribe(
      this.renderKube.bind(this),
      catchError((error) => of(error)),
    );
  }

  private initDialog(target) {
    const popupWidth = 250;
    const dialogRef = this.dialog.open(ConfirmModalComponent, {
      width: `${popupWidth}px`,
    });
    dialogRef.updatePosition({
      top: `${target.clientY}px`,
      left: `${target.clientX - popupWidth - 10}px`,
    });
    return dialogRef;
  }
}

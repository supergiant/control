import { of, Subscription, timer as observableTimer } from 'rxjs';
import { catchError, filter, switchMap } from 'rxjs/operators';
import { ChangeDetectionStrategy, Component, OnDestroy, OnInit, ViewChild, ViewEncapsulation, Inject } from '@angular/core';
import { Location } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { MatDialog, MatPaginator, MatSort, MatTableDataSource } from '@angular/material';
import { animate, state, style, transition, trigger } from '@angular/animations';

import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { UtilService } from '../../shared/supergiant/util/util.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { ConfirmModalComponent } from '../../shared/modals/confirm-modal/confirm-modal.component';
import { DeleteClusterModalComponent } from './delete-cluster-modal/delete-cluster-modal.component';
import { DeleteReleaseModalComponent } from './delete-release-modal/delete-release-modal.component';
import { SshCommandsModalComponent } from './ssh-commands-modal/ssh-commands-modal.component';
import { KubectlConfigModalComponent } from './kubectl-config-modal/kubectl-config-modal.component';
import { TaskLogsComponent } from './task-logs/task-logs.component';
import { ReleaseInfoModalComponent } from './release-info-modal/release-info-modal.component';
import { WINDOW } from '../../shared/helpers/window-providers';


@Component({
  selector: 'app-cluster',
  templateUrl: './cluster.component.html',
  styleUrls: [ './cluster.component.scss' ],
  changeDetection: ChangeDetectionStrategy.Default,
  encapsulation: ViewEncapsulation.None,
  animations: [
      trigger('detailExpand', [
        state('collapsed', style({height: '0px', minHeight: '0', visibility: 'hidden'})),
        state('expanded', style({height: '*', visibility: 'visible'})),
        transition('expanded <=> collapsed', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
      ]),
    ]
})

export class ClusterComponent implements OnInit, OnDestroy {
  clusterId: number;
  subscriptions = new Subscription();
  public kube: any;
  public kubeString: string;

  // machine list vars
  machines: any;
  machineListColumns = ["state", "role", "size", "name", "cpu", "ram", "region", "publicIp", "delete"];

  // task list vars
  tasks: any;
  taskListColumns = ["status", "type", "id", "steps", "logs"];
  expandedTaskIds = new Set();

  releases: any;
  releaseListColumns = ["status", "name", "chart", "chartVersion", "version", "lastDeployed", "info", "delete"];

  masterTasksStatus = "executing";
  nodeTasksStatus = "queued";
  clusterTasksStatus = "queued";

  cpuUsage: number;
  ramUsage: number;
  machineMetrics = {};

  kubectlConfig: any;

  clusterServices: any
  serviceListColumns = ["name", "type", "namespace", "selfLink"];


  constructor(
    private route: ActivatedRoute,
    private location: Location,
    private router: Router,
    private supergiant: Supergiant,
    private util: UtilService,
    private notifications: Notifications,
    public dialog: MatDialog,
    public http: HttpClient,
    @Inject(WINDOW) private window: Window
  ) {
      route.params.subscribe(params => {
        this.clusterId = params.id;
        this.getKube();
      })
    }

  @ViewChild(MatSort) sort: MatSort;
  @ViewChild(MatPaginator) paginator: MatPaginator;

  ngOnInit() {
    this.clusterId = this.route.snapshot.params.id;
    this.getKube();
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  combineAndFlatten(objOne, objTwo) {
    const arr = [];

    Object.keys(objOne).forEach((key) => {
      arr.push(objOne[key]);
    });

    Object.keys(objTwo).forEach((key) => {
      arr.push(objTwo[key]);
    });

    return arr;
  }

  getKubeStatus(clusterId) {
    // we should make Tasks a part of the Supergiant instance
    // if we start using them outside of this
    return this.util.fetch("/v1/api/kubes/" + clusterId + "/tasks");
  }

  toggleSteps(task) {

    task.showSteps = !task.showSteps;

    if (this.expandedTaskIds.has(task.id)) {
      this.expandedTaskIds.delete(task.id);
    } else { this.expandedTaskIds.add(task.id); }
  }

  taskComplete(task) {
    return task.stepsStatuses.every((s) => s.status == "success");
  }

  restartTask(taskId) {
    this.tasks.data.map(t => {
      if (t.id == taskId) {
        t.restarting = true;
      }}
    );

    this.util.post("v1/api/tasks/" + taskId + "/restart", {}).subscribe(
      res => console.log(res),
      err => console.error(err)
    );
  }

  viewTaskLog(taskId) {
    const modal = this.dialog.open(TaskLogsComponent, {
      width: "1080px",
      data: { taskId: taskId, hostname: this.window.location.hostname }
    })
  }

  setProvisioningStep(tasks) {
    const masterPatt = /master/g;
    const masterTasks = tasks.filter(t => {
      return masterPatt.test(t.type.toLowerCase());
    });

    const nodePatt = /node/g;
    const nodeTasks = tasks.filter(t => {
      return nodePatt.test(t.type.toLowerCase());
    });

    // oh my god I'm so sorry
    if (masterTasks.every(this.taskComplete)) {
      // masters complete
      this.masterTasksStatus = "complete";

      if (nodeTasks.every(this.taskComplete)) {
        // masters AND nodes complete
        this.nodeTasksStatus = "complete";
        this.clusterTasksStatus = "executing";

      } else {
        // masters complete, nodes executing
        this.nodeTasksStatus = "executing"
      }
    } else {
      // masters executing
      this.masterTasksStatus = "executing";
    }
  }

  getKube() {
    // TODO: shameful how smart this ENTIRE component has become.
    this.subscriptions.add(observableTimer(0, 10000).pipe(
      switchMap(() => this.supergiant.Kubes.get(this.clusterId))).subscribe(
        k => {
          this.kube = k;
          // for dev-ing
          // this.kube.state = "provisioning";

          switch (this.kube.state) {
            case "operational": {
              this.renderMachines(this.kube);
              this.getReleases();
              this.getClusterMetrics();
              this.getMachineMetrics();
              this.getKubectlConfig();
              this.getClusterServices();
              break;
            }
            case "provisioning": {
              this.getKubeStatus(this.clusterId).subscribe(
                tasks => {
                  this.setProvisioningStep(tasks);

                  const rows = [];
                  tasks.forEach(t => {
                    if (this.expandedTaskIds.has(t.id)) {
                      t.showSteps = true;
                    };
                    rows.push(t, { detailRow: true, t })
                  });
                  this.tasks = new MatTableDataSource(rows);
                  this.tasks.sort = this.sort;
                  this.tasks.paginator = this.paginator;
                },
                err => console.log(err)
              )
              break;
            }
            case "failed": {
              this.getKubeStatus(this.clusterId).subscribe(
                tasks => {

                  const rows = [];
                  tasks.forEach(t => {
                    if (this.expandedTaskIds.has(t.id)) {
                      t.showSteps = true;
                    };
                    rows.push(t, { detailRow: true, t })
                  });
                  this.tasks = new MatTableDataSource(rows);
                  this.tasks.sort = this.sort;
                  this.tasks.paginator = this.paginator;
                },
                err => console.log(err)
              );
              this.masterTasksStatus = "failed";
              this.nodeTasksStatus = "failed";
              this.clusterTasksStatus = "failed";
              break;
            }
            default: {
              break;
            }
          }
        },
        err => console.error(err)
      ))
  }

  renderMachines(kube) {
    const machineMetrics = {};
    const masterNames = Object.keys(kube.masters);
    const nodeNames = Object.keys(kube.nodes);

    masterNames.forEach(name => {
      const lowercaseName = name.toLowerCase();
      if (this.machineMetrics[lowercaseName]) {
        kube.masters[name].metrics = this.machineMetrics[lowercaseName];
      }
    })

    nodeNames.forEach(name => {
      const lowercaseName = name.toLowerCase();
      if (this.machineMetrics[lowercaseName]) {
        kube.nodes[name].metrics = this.machineMetrics[lowercaseName];
      }
    })

    this.machines = new MatTableDataSource(this.combineAndFlatten(kube.masters, kube.nodes));
    this.machines.sort = this.sort;
    this.machines.paginator = this.paginator;
  }

  getReleases() {
    this.supergiant.HelmReleases.get(this.clusterId).subscribe(
      res => {
        const releases = res.filter(r => r.status != "DELETED")
        this.releases = new MatTableDataSource(releases);
      },
      err => console.error(err)
    )
  }

  getClusterMetrics() {
    this.supergiant.Kubes.getClusterMetrics(this.clusterId).subscribe(
      res => {
        this.cpuUsage = res.cpu;
        this.ramUsage = res.memory;
      },
      err => console.error(err)
    )
  }

  getMachineMetrics() {
    this.supergiant.Kubes.getMachineMetrics(this.clusterId).subscribe(
      res => {
        this.machineMetrics = this.calculateMachineMetrics(res);
        this.renderMachines(this.kube)
      },
      err => console.error(err)
    )
  }

  getKubectlConfig() {
    // TODO: move to service
    this.util.fetch('v1/api/kubes/' + this.clusterId + '/users/kubernetes-admin/kubeconfig').subscribe(
      res => this.kubectlConfig = res,
      err => console.error(err)
    )
  }

  getClusterServices() {
    this.supergiant.Kubes.getClusterServices(this.clusterId).subscribe(
      res => this.clusterServices = new MatTableDataSource(res),
      err => console.error(err)
    )
  }

  calculateMachineMetrics(machines) {
    Object.keys(machines).forEach(m => {
      machines[m].cpu = (machines[m].cpu * 100).toFixed(1);
      machines[m].memory = (machines[m].memory * 100).toFixed(1);
    });

    return machines;
  }

  removeNode(nodeName: string, target) {
    const dialogRef = this.initDialog(target);

    dialogRef
      .afterClosed()
      .pipe(
        filter(isConfirmed => isConfirmed),
        switchMap(() => this.supergiant.Nodes.delete(this.clusterId, nodeName)),
        switchMap(() => this.supergiant.Kubes.get(this.clusterId)),
        catchError((error) => of(error)),
      ).subscribe(
        k => this.renderMachines(k),
        err => {
          console.error(err);
          this.error(this.clusterId, err)
        },
     );
  }

  deleteCluster() {
    const dialogRef = this.initDeleteCluster(this.kube.state);

    dialogRef
      .afterClosed()
      .pipe(
        filter(isConfirmed => isConfirmed),
        switchMap(() => this.supergiant.Kubes.delete(this.clusterId)),
        catchError((error) => of(error)),
      ).subscribe(
        res => {
          this.success(this.clusterId);
          this.router.navigate([""]);
        },
        err => {
          console.error(err);
          this.error(this.clusterId, err)
        },
     );
  }

  deleteRelease(releaseName, event) {
    const dialogRef = this.initDeleteRelease(releaseName)

    dialogRef
      .afterClosed()
      .pipe(
        filter(res => res.deleteRelease),
        switchMap(res => this.supergiant.HelmReleases.delete(releaseName, this.clusterId, !res.deleteConfigs)),
        catchError(err => of(err))
      ).subscribe(
        res => this.getReleases(),
        err => console.error(err)
      )
  }

  showSshCommands() {
    const masters = [];
    const nodes = [];

    Object.keys(this.kube.masters).map(m => masters.push(this.kube.masters[m]));
    Object.keys(this.kube.nodes).map(m => nodes.push(this.kube.nodes[m]));

    this.initSshCommands(masters, nodes);
  }

  showKubectlConfig() {
    this.initKubectlConfig(this.kubectlConfig);
  }

  showReleaseInfo(releaseName) {
    this.initReleaseInfo(releaseName);
  }

  openService(proxyPort) {
    const hostname = this.window.location.hostname;
    const link = "http://" + hostname + ":" + proxyPort;

    this.window.open(link);
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

  private initDeleteCluster(clusterState) {
    const dialogRef = this.dialog.open(DeleteClusterModalComponent, {
      width: "500px",
      data: { state: clusterState }
    });
    return dialogRef;
  }

  private initDeleteRelease(name) {
    const dialogRef = this.dialog.open(DeleteReleaseModalComponent, {
      width: "max-content",
      data: { name: name }
    })
    return dialogRef;
  }

  private initSshCommands(masters, nodes) {
    const dialogRef = this.dialog.open(SshCommandsModalComponent, {
      width: "600px",
      data: { masters: masters, nodes: nodes }
    })

    return dialogRef;
  }

  private initKubectlConfig(config) {
    const dialogRef = this.dialog.open(KubectlConfigModalComponent, {
      width: "800px",
      data: { config: config }
    })

    return dialogRef;
  }

  private initReleaseInfo(releaseName) {
    const dialogRef = this.dialog.open(ReleaseInfoModalComponent, {
      width: "800px",
      data: { releaseName: releaseName, clusterId: this.clusterId }
    })

    return dialogRef;
  }

  expandRow = (_, row) => row.hasOwnProperty('detailRow');

  success(name) {
    this.notifications.display(
      'success',
      'Kube: ' + name,
      'Deleted...',
    );
  }

  error(name, error) {
    this.notifications.display(
      'error',
      'Kube: ' + name,
      'Error:' + error.statusText);
  }
}

import { of, Subscription, timer as observableTimer } from 'rxjs';
import { catchError, filter, switchMap } from 'rxjs/operators';
import { ChangeDetectionStrategy, Component, OnDestroy, OnInit, ViewChild, ViewEncapsulation } from '@angular/core';
import { Location } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { SafeResourceUrl } from '@angular/platform-browser';
import { MatDialog, MatPaginator, MatSort, MatTableDataSource } from '@angular/material';
import {animate, state, style, transition, trigger} from '@angular/animations';

import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { UtilService } from '../../shared/supergiant/util/util.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { ConfirmModalComponent } from '../../shared/modals/confirm-modal/confirm-modal.component';
import { TaskLogsComponent } from './task-logs/task-logs.component';


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
  id: number;
  subscriptions = new Subscription();
  public kube: any;
  public kubeString: string;

  // machine list vars
  machines: any;
  machineListColumns = ["state", "role", "name", "id", "region", "publicIp"];

  // task list vars
  tasks: any;
  taskListColumns = ["status", "type", "id", "steps", "logs"];
  expandedTaskId: any;

  // temp for demo, remove ASAP
  mastersStatus = "executing";
  nodesStatus = "queued";
  readyStatus = "queued";

  constructor(
    private route: ActivatedRoute,
    private location: Location,
    private router: Router,
    private supergiant: Supergiant,
    private util: UtilService,
    private notifications: Notifications,
    public dialog: MatDialog,
    public http: HttpClient,
  ) { }

  @ViewChild(MatSort) sort: MatSort;
  @ViewChild(MatPaginator) paginator: MatPaginator;

  ngOnInit() {
    this.id = this.route.snapshot.params.id;
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

  getKubeStatus(id) {
    // we should make Tasks a part of the Supergiant instance
    // if we start using them outside of this
    return this.util.fetch("/v1/api/kubes/" + id + "/tasks");
  }

  toggleSteps(task) {

    task.showSteps = !task.showSteps;

    if (task.id != this.expandedTaskId) {
      this.expandedTaskId = task.id;
    } else { this.expandedTaskId = undefined; }
  }

  taskComplete(task) {
    // return task.status == "complete";
    return task.stepsStatuses.every((s) => s.status == "success");
  }

  restartTask(taskId) {
    this.tasks.data.map(t => {
      if (t.id == taskId) {
        t.restarting = true;
      }}
    );

    this.util.post("/tasks/" + taskId + "/restart", {}).subscribe(
      res => console.log(res),
      err => console.error(err)
    );
  }

  viewTaskLog(taskId) {
    const modal = this.dialog.open(TaskLogsComponent, {
      width: "1080px",
      data: { taskId: taskId }
    })
  }

  setProvisioningStep(tasks) {
    const masterPatt = /master/g;
    const masterTasks = tasks.filter(t => {
      return masterPatt.test(t.type.toLowerCase());
    })

    const nodePatt = /node/g;
    const nodeTasks = tasks.filter(t => {
      return nodePatt.test(t.type.toLowerCase());
    })

    // oh my god I'm so sorry
    if (masterTasks.every(this.taskComplete)) {
      // masters complete
      this.mastersStatus = "complete";

      if (nodeTasks.every(this.taskComplete)) {
        // masters AND nodes complete
        this.nodesStatus = "complete";
        this.readyStatus = "executing";

        // TODO: what is technically "ready?"
        setTimeout(() => { this.readyStatus = "complete"; console.log("ready") }, 5000);
      } else {
        // masters complete, nodes executing
        this.nodesStatus = "executing"
      }
    } else {
      // masters executing
      this.mastersStatus = "executing";
    }
  }

  getKube() {
    // TODO: shameful how smart this ENTIRE component has become.
    this.subscriptions.add(observableTimer(0, 10000).pipe(
      switchMap(() => this.supergiant.Kubes.get(this.id))).subscribe(
        k => {
          this.kube = k;
          // for dev-ing
          // this.kube.state = "provisioning";

          switch (k.state) {
            case "operational": {
              this.renderKube(this.kube);
              break;
            }
            case "provisioning": {
              this.getKubeStatus(this.id).subscribe(
                tasks => {
                  // for demo:
                  // do some task overview checking here
                  // e.g. are we on masters? nodes? set vars accordingly
                  this.setProvisioningStep(tasks);

                  const rows = [];
                  tasks.forEach(t => {
                    if (t.id == this.expandedTaskId) {
                      t.showSteps = true;
                    };
                    // t.status = "complete";
                    // t.status = "failed";
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
              console.log("deployment has failed")
              // duped from 'provisioning' case for demo
              this.getKubeStatus(this.id).subscribe(
                tasks => {

                  const rows = [];
                  tasks.forEach(t => {
                    if (t.id == this.expandedTaskId) {
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
              this.mastersStatus = "failed";
              this.nodesStatus = "failed";
              this.readyStatus = "failed";
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

  renderKube(kube) {

    let masters = kube.masters;
    let nodes = kube.nodes;

    this.machines = new MatTableDataSource(this.combineAndFlatten(masters, nodes));
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
        k => this.renderKube(k),
        err => console.error(err),
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

  expandRow = (_, row) => row.hasOwnProperty('detailRow');
}

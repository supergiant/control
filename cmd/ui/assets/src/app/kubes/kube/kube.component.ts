import { Component, Input, AfterViewInit } from '@angular/core';
import { KubesService } from '../kubes.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { NgbProgressbarConfig } from '@ng-bootstrap/ng-bootstrap';


@Component({
  selector: '[app-kube]',
  templateUrl: './kube.component.html',
  styleUrls: ['./kube.component.css'],
})
export class KubeComponent implements AfterViewInit {
  @Input() kube: any;
  private progress = false;
  private progressValue: number;
  constructor(
    private kubesService: KubesService,
    private notifications: Notifications,
    private config: NgbProgressbarConfig,
  ) {
    // config.max = 1000;
    config.striped = true;
    config.animated = true;
    // config.type = 'success';
  }

  ngAfterViewInit() {
    setTimeout(_ => this.status(this.kube)); // BUGFIX: https://github.com/angular/angular/issues/6005#issuecomment-165911194
  }

  status(kube) {
    if (kube.status && kube.status.steps_completed && kube.status.total_steps) {
      this.progress = true;
      this.config.max = kube.status.total_steps;
      this.progressValue = kube.status.steps_completed;

      switch (kube.status.description) {
        case 'deleting':
          this.config.type = 'danger';
          break;
        case 'provisioning':
          this.config.type = 'success';
          break;
        case 'updating':
          this.config.type = 'warning';
          break;
      }
    }



    if (kube.status && kube.status.error && kube.status.retries === kube.status.max_retries) {
      return 'status status-danger';
    } else if (kube.status) {
      return 'status status-transitioning';
    } else if (kube.passive_status && !kube.passive_status_okay) {
      return 'status status-warning';
    } else {
      return 'status status-ok';
    }
  }
}

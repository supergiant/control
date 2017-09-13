import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { ActivatedRoute, Router, Params } from '@angular/router';
import { Observable } from 'rxjs/Observable';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { SystemModalService } from '../../shared/system-modal/system-modal.service';
import { LoginComponent } from '../../login/login.component';

@Component({
  selector: 'app-app-details',
  templateUrl: './app-details.component.html',
  styleUrls: ['./app-details.component.css']
})
export class AppDetailsComponent implements OnInit, OnDestroy {

  id: number;
  subscriptions = new Subscription();
  app: any;
  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private supergiant: Supergiant,
    private notifications: Notifications,
    private systemModalService: SystemModalService,
    public loginComponent: LoginComponent,
  ) { }

  ngOnInit() {
    this.id = this.route.snapshot.params.id;
    this.getAccount();
  }

  openSystemModal(message) {
    this.systemModalService.openSystemModal(message);
  }

  getAccount() {
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.HelmCharts.get(this.id)).subscribe(
      (chart) => { this.app = chart; },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  goBack() {
    this.router.navigate(['/apps']);
  }
  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

}

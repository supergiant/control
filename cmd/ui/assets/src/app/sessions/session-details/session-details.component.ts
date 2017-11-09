import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { ActivatedRoute, Router } from '@angular/router';
import { Observable } from 'rxjs/Observable';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { SystemModalService } from '../../shared/system-modal/system-modal.service';
import { LoginComponent } from '../../login/login.component';

@Component({
  selector: 'app-session-details',
  templateUrl: './session-details.component.html',
  styleUrls: ['./session-details.component.scss']
})
export class SessionDetailsComponent implements OnInit, OnDestroy {
  id: number;
  subscriptions = new Subscription();
  session: any;
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
    this.getSession();
  }

  openSystemModal(message) {
    this.systemModalService.openSystemModal(message);
  }

  getSession() {
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.Sessions.get(this.id)).subscribe(
      (session) => { this.session = session; },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  goBack() {
    this.router.navigate(['/sessions']);
  }
  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

}

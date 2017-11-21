import { Component, OnDestroy } from '@angular/core';
import { Supergiant } from '../shared/supergiant/supergiant.service';
import { CookieMonster } from '../shared/cookies/cookies.service';
import { Router } from '@angular/router';
import { Observable } from 'rxjs/Observable';
import { Subscription } from 'rxjs/Subscription';
import { SessionModel } from './session.model';
import { Notifications } from '../shared/notifications/notifications.service';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnDestroy {
  public username: string;
  public password: string;
  private session: any;
  private id: string;
  private sessionCookie: string;
  private previousUrl: string;
  private refresh: boolean;
  private subscriptions = new Subscription();
  public status: string;

  constructor(
    private supergiant: Supergiant,
    private router: Router,
    private cookieMonster: CookieMonster,
    private notifications: Notifications,
  ) { }

  validateUser() {
    this.sessionCookie = this.cookieMonster.getCookie('session');
    if (this.sessionCookie) {
      this.supergiant.UtilService.sessionToken = 'SGAPI session="' + this.sessionCookie + '"';
      this.supergiant.sessionID = this.sessionCookie;
    }

    return this.supergiant.Sessions.valid(this.supergiant.sessionID);
  }

  handleError() {
    return Observable.of(false);
  }

  error(msg) {
    this.notifications.display(
      'error',
      'Login Error',
      'Error:' + msg);
  }

  onSubmit() {
    this.status = 'status status-transitioning';
    const creds = new SessionModel;
    creds.session.model.user.username = this.username;
    creds.session.model.user.password = this.password;

    this.subscriptions.add(this.supergiant.Sessions.create(creds.session.model).subscribe(
      (session) => {
        this.session = session;
        this.supergiant.UtilService.sessionToken = 'SGAPI session="' + this.session.id + '"';
        this.supergiant.sessionID = this.session.id;
        this.cookieMonster.setCookie({ name: 'session', value: this.session.id, secure: true });

        const countdown = Observable
          .interval(100)
          .take(50) // 5 seconds
          .subscribe(y => {
            if (this.cookieMonster.getCookie('session') === this.session.id) {
              this.supergiant.loginSuccess = true;
              this.router.navigate(['/dashboard']);
              countdown.unsubscribe();
            }

            if (!this.supergiant.loginSuccess && y === 49) {
              this.status = 'status status-danger';
              this.error('No Login Cookie Found');
            }
          });

      },
      (err) => {
        this.status = 'status status-danger';
        this.error('Invalid Login');
      }
    ));
  }

  logOut() {
    this.subscriptions.add(this.supergiant.Sessions.delete(this.supergiant.sessionID).subscribe(
      (session) => {
        this.supergiant.sessionID = '';
        this.cookieMonster.deleteCookie('session');
        this.supergiant.loginSuccess = false;
        this.router.navigate(['/login']);
      }
    ));
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

}

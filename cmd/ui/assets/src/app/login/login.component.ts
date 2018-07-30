import { Component, OnDestroy } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Supergiant } from '../shared/supergiant/supergiant.service';
import { Router } from '@angular/router';
import { Observable } from 'rxjs/Observable';
import { Subscription } from 'rxjs/Subscription';
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
    private notifications: Notifications,
  ) { }

  error(msg) {
    this.notifications.display(
      'error',
      'Login Error',
      'Error:' + msg);
  }

  onSubmit() {
    this.status = 'status status-transitioning';
    const creds = { "login": this.username, "password": this.password };

    this.supergiant.Auth.login(creds).subscribe(
      res => {
        if (res["status"] === 200) {
          this.supergiant.loginSuccess = true;
          this.router.navigate(['/dashboard']);
        }
      },
      err => {
        this.status = 'status status-danger';
        this.error('Invalid Login');
      }
    )
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

}

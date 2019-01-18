import { Component, OnDestroy, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { FormBuilder, FormGroup, Validators, ValidatorFn, AbstractControl } from '@angular/forms';
import { Observable, Subscription, of } from 'rxjs';
import { catchError, switchMap, zip } from 'rxjs/operators';

import { Supergiant } from '../shared/supergiant/supergiant.service';
import { UtilService } from '../shared/supergiant/util/util.service';
import { Notifications } from '../shared/notifications/notifications.service';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnDestroy, OnInit {
  private subscriptions = new Subscription();
  public isColdStart: boolean;
  public loginForm: FormGroup;

  constructor(
    private supergiant: Supergiant,
    private router: Router,
    private notifications: Notifications,
    private util: UtilService,
    private formBuilder: FormBuilder
  ) { }

  initColdStartForm() {
    const pattern = new RegExp('(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])[a-zA-Z0-9]')
    this.loginForm.controls['password'].setValidators([
      Validators.required,
      Validators.minLength(8),
      Validators.pattern(pattern)
    ]);
    this.loginForm.controls['password'].updateValueAndValidity()
  }

  error(msg) {
    this.notifications.display(
      'error',
      'Login Error',
      'Error:' + msg);
  }

  login(creds) {
    this.supergiant.Auth.login(creds).subscribe(
      res => {
        this.supergiant.loginSuccess = true;
        this.router.navigate(['/dashboard']);
      },
      err => {
        this.error(err.userMessage)
      }
    )
  }

  onSubmit() {
    const creds = this.loginForm.value

    if (this.isColdStart) {
      this.supergiant.Auth.signup(creds).pipe(
        switchMap(res => this.supergiant.Auth.login(creds)),
        catchError((err) => of(err))
      ).subscribe(
        res => {
          this.supergiant.loginSuccess = true;
          this.router.navigate(['/dashboard']);
        },
        err => this.error(err.userMessage)
      )
    } else {
      this.login(creds)
    }
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  ngOnInit() {
    this.loginForm = this.formBuilder.group({
      login: ['', Validators.required],
      password: ['', Validators.required]
    })

    this.util.fetch('/coldstart').subscribe(
      res => {
        this.isColdStart = res.isColdStart;

        if (this.isColdStart) {
          this.initColdStartForm();
        }
      },
      err => this.error(err.userMessage)
    )
  }

}

import { Component, OnInit } from '@angular/core';
import { Supergiant } from '../shared/supergiant/supergiant.service'
import { CookieMonster } from '../shared/cookies/cookies.service'
import { Router } from '@angular/router';
import { Observable } from 'rxjs/Rx';
import 'rxjs/add/operator/catch';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {
  private username: string;
  private password: string;
  private session: any;
  private id: string;
  private sessionCookie: string;
  previousUrl: string;
  private refresh:boolean;

  constructor(
    private supergiant: Supergiant,
    private router: Router,
    private cookieMonster: CookieMonster,
  ) { }

  ngOnInit() {}

  validateUser() {
    this.sessionCookie = this.cookieMonster.getCookie('session')
    if (this.sessionCookie) {
      this.supergiant.UtilService.sessionToken = 'SGAPI session="'+ this.sessionCookie +'"'
      this.supergiant.sessionID = this.sessionCookie
    }

    return this.supergiant.Sessions.valid(this.supergiant.sessionID)
  }

  handleError() {
    return Observable.of(false);
  }

  onSubmit() {
    let creds = '{"user":{"username":"'+ this.username +'", "password":"'+ this.password +'"}}'
    this.supergiant.Sessions.create(JSON.parse(creds)).subscribe(
      (session) => { this.session = session
        this.supergiant.UtilService.sessionToken = 'SGAPI session="'+ this.session.id +'"'
        this.supergiant.sessionID = this.session.id
        this.cookieMonster.setCookie({name:'session',value:this.session.id, secure:true });
        this.supergiant.loginSuccess = true
        this.router.navigate(['/kubes']);
      }
    )
  }

  logOut() {
    this.supergiant.Sessions.delete(this.supergiant.sessionID).subscribe(
      (session) => {
        console.log(session)
        this.supergiant.sessionID = ''
        this.cookieMonster.deleteCookie('session')
        this.supergiant.loginSuccess = false
        this.router.navigate(['/login']);
      }
    )
  }

}

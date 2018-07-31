import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../shared/supergiant/auth/auth.service';

@Component({
  selector: 'app-navigation',
  templateUrl: './navigation.component.html',
  styleUrls: ['./navigation.component.scss']
})
export class NavigationComponent {
  public isCollapsed: boolean;
  constructor(
    public auth: AuthService,
    public router: Router,
  ) { }

  logout() {
  	this.auth.logout();
  	this.router.navigate(['']);
  }
}

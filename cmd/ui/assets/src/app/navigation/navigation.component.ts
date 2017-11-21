import { Component } from '@angular/core';
import { LoginComponent } from '../login/login.component';

@Component({
  selector: 'app-navigation',
  templateUrl: './navigation.component.html',
  styleUrls: ['./navigation.component.scss']
})
export class NavigationComponent {
  public isCollapsed: boolean;
  constructor(
    public loginComponent: LoginComponent,
  ) { }
}

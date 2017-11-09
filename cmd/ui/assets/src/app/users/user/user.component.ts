import { Component, Input } from '@angular/core';
import { UsersService } from '../users.service';

@Component({
  selector: '[app-user]', // tslint:disable-line
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.scss']
})
export class UserComponent {
  @Input() user: any;
  constructor(public usersService: UsersService) { }
}

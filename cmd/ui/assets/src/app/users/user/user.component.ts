import { Component, Input} from '@angular/core';
import { UsersService } from '../users.service';

@Component({
  selector: '[app-user]',
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.css']
})
export class UserComponent {
  @Input() user: any;
  constructor(private usersService: UsersService) { }
}

import { Component, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material';

import { MenuModalComponent } from './menu-modal/menu-modal.component';

@Component({
  selector: 'app-user-menu',
  templateUrl: './user-menu.component.html',
  styleUrls: ['./user-menu.component.scss']
})
export class UserMenuComponent implements OnInit {

  constructor( public dialog: MatDialog ) { }

  toggleMenu(event) {
    const menu = this.initDialog(event)
  }

  initDialog(event) {
    const popupWidth = 200;
    const dialogRef = this.dialog.open(MenuModalComponent, {
      width: `${popupWidth}px`,
      backdropClass: "backdrop"
    });
    dialogRef.updatePosition({
      top: `${event.clientY + 20}px`,
      left: `${event.clientX - popupWidth}px`,
    });
    return dialogRef;
  }

  ngOnInit() {
  }

}

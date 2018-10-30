import { Component, OnInit } from '@angular/core';
import { MatDialogRef } from '@angular/material';
import { Router } from '@angular/router'

import { AuthService } from '../../../shared/supergiant/auth/auth.service'

@Component({
  selector: 'menu-modal',
  templateUrl: './menu-modal.component.html',
  styleUrls: ['./menu-modal.component.scss']
})
export class MenuModalComponent implements OnInit {

  constructor(
    public dialogRef: MatDialogRef<MenuModalComponent>,
    private router: Router,
    private auth: AuthService
  ) { }

  navigate(path) {
    this.router.navigate([path]);
    this.dialogRef.close();
  }

  logout() {
    this.auth.logout();
    this.dialogRef.close();
    this.router.navigate([""])
  }

  ngOnInit() {
  }

}

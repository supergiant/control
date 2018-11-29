import { Component, OnInit, Inject } from '@angular/core';
import { MAT_DIALOG_DATA } from '@angular/material';

@Component({
  selector: 'ssh-commands-modal',
  templateUrl: './ssh-commands-modal.component.html',
  styleUrls: ['./ssh-commands-modal.component.scss']
})
export class SshCommandsModalComponent implements OnInit {

  constructor(@Inject(MAT_DIALOG_DATA) public data: any) { }

  ngOnInit() {
  }

}

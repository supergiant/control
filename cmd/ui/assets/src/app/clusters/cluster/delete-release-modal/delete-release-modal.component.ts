import { Component, OnInit, Inject, ViewChild } from '@angular/core';
import { MAT_DIALOG_DATA } from '@angular/material';

@Component({
  selector: 'delete-release-modal',
  templateUrl: './delete-release-modal.component.html',
  styleUrls: ['./delete-release-modal.component.scss']
})
export class DeleteReleaseModalComponent implements OnInit {

  constructor(@Inject(MAT_DIALOG_DATA) public data: any) { }

  // @ViewChild("keepConfig") config;

  ngOnInit() {
  }

}

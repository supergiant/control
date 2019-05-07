import { Component, OnInit, Inject } from '@angular/core';
import { MAT_DIALOG_DATA } from '@angular/material';

@Component({
  selector: 'public-key-modal',
  templateUrl: './public-key-modal.component.html',
  styleUrls: ['./public-key-modal.component.scss']
})
export class PublicKeyModalComponent implements OnInit {

  constructor(@Inject(MAT_DIALOG_DATA) public data: any) { }

  key: string;

  ngOnInit() {
  	this.key = this.data.key;
  }

}

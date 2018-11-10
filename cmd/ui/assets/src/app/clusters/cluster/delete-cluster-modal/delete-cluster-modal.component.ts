import { Component, OnInit, Inject } from '@angular/core';
import { MAT_DIALOG_DATA } from '@angular/material';


@Component({
  selector: 'delete-cluster-modal',
  templateUrl: './delete-cluster-modal.component.html',
  styleUrls: ['./delete-cluster-modal.component.scss']
})
export class DeleteClusterModalComponent implements OnInit {

  constructor(@Inject(MAT_DIALOG_DATA) public data: any) { }

  ngOnInit() {
  }

}

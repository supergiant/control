import { Component, OnInit, Inject } from '@angular/core';
import { MAT_DIALOG_DATA } from '@angular/material';

@Component({
  selector: 'kubectl-config-modal',
  templateUrl: './kubectl-config-modal.component.html',
  styleUrls: ['./kubectl-config-modal.component.scss']
})
export class KubectlConfigModalComponent implements OnInit {

  constructor(@Inject(MAT_DIALOG_DATA) public data: any) { }

  config: string;

  ngOnInit() {
    this.config = JSON.stringify(this.data.config, null, 2);
  }

}

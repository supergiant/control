import { Component, Inject, OnInit }     from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material";

@Component({
  selector: 'configure',
  templateUrl: './configure.component.html',
  styleUrls: ['./configure.component.scss']
})
export class ConfigureComponent implements OnInit {
  values: string;

  constructor(
    private dialogRef: MatDialogRef<ConfigureComponent>,
    @Inject(MAT_DIALOG_DATA) public dialogData
  ) {
  }

  ngOnInit() {
    this.dialogData.chart$
      .subscribe(chart => this.values = chart.values);
  }

  onChange(values) {
    this.values = values;
  }

  save() {
    this.dialogRef.close(this.values);
  }

}

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
    this.dialogData.values
      .subscribe(values => this.values = values);
  }

  onChange(e) {
    console.log(e);
  }

  save() {
    this.dialogRef.close(this.values);
  }

}

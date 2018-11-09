import { Component }    from '@angular/core';
import { MatDialogRef } from "@angular/material";

@Component({
  selector: 'remove-repo-dialog',
  templateUrl: './remove-repo-dialog.component.html',
  styleUrls: ['./remove-repo-dialog.component.scss']
})
export class RemoveRepoDialogComponent {

  constructor(
    private dialogRef: MatDialogRef<RemoveRepoDialogComponent>,
  ) {
  }

  confirm() {
    this.dialogRef.close(true);
  }
}

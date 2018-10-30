import { Component, OnInit, Inject } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';

@Component({
  selector: 'cluster-list-modal',
  templateUrl: './cluster-list-modal.component.html',
  styleUrls: ['./cluster-list-modal.component.scss'],
})
export class ClusterListModalComponent implements OnInit {

  constructor(
    private dialogRef: MatDialogRef<ClusterListModalComponent>,
    private router: Router,
    private route: ActivatedRoute,
    @Inject(MAT_DIALOG_DATA) private data: any
  ) { this.clusters = this.data.clusters }

  public clusters: Array<any>;

  navigate(name) {
    this.router.navigate(['/clusters/', name]);
    this.dialogRef.close();
  }

  newCluster() {
    this.router.navigate(['/clusters/new']);
    this.dialogRef.close();
  }

  ngOnInit() {
  }

}

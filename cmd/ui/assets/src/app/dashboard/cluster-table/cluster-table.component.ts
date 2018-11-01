import { Component, OnInit, Input } from '@angular/core';
import { MatTableDataSource } from '@angular/material';
import { Router } from '@angular/router';

@Component({
  selector: 'cluster-table',
  templateUrl: './cluster-table.component.html',
  styleUrls: ['./cluster-table.component.scss']
})
export class ClusterTableComponent implements OnInit {

  public clusterColumns = ["state", "accountName", "region", "k8sversion", "mastersCount", "nodesCount", "operatingSystem", "dockerVersion", "helmVersion"];

  constructor( private router: Router ) { }

  @Input() cluster: any

  ngOnInit() {
    console.log(this.cluster);
  }

}

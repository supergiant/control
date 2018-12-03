import { Component, OnInit, Inject } from '@angular/core';
import { MAT_DIALOG_DATA } from '@angular/material';

import { Supergiant } from '../../../shared/supergiant/supergiant.service';

@Component({
  selector: 'release-info-modal',
  templateUrl: './release-info-modal.component.html',
  styleUrls: ['./release-info-modal.component.scss']
})
export class ReleaseInfoModalComponent implements OnInit {

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: any,
    private supergiant: Supergiant
  ) { }

  clusterId: string;
  releaseName: string;
  releaseInfo: string;

  getReleaseInfo() {
    this.supergiant.HelmReleases.get(this.clusterId, this.releaseName).subscribe(
      res => this.releaseInfo = res,
      err => console.error(err)
    )
  }

  ngOnInit() {
    this.releaseName = this.data.releaseName;
    this.clusterId = this.data.clusterId;
    this.getReleaseInfo();
  }

}

import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';

import { Supergiant } from '../../../shared/supergiant/supergiant.service'

@Component({
  selector: 'service-proxy',
  templateUrl: './service-proxy.component.html',
  styleUrls: ['./service-proxy.component.scss']
})
export class ServiceProxyComponent implements OnInit {

  constructor(
    private supergiant: Supergiant,
    private route: ActivatedRoute,
    private router: Router
  ) { }

  clusterId: number;
  serviceUrl: any;
  body: any;

  getServiceUi() {
    this.supergiant.Kubes.getClusterServiceUi(this.clusterId, this.serviceUrl).subscribe(
      res => {
        console.log(res);
        this.body = res;
      },
      err => console.error(err)
    )
  }

  ngOnInit() {
    this.clusterId = this.route.snapshot.params.id;
    this.serviceUrl = this.route.snapshot.queryParams.url;
    this.getServiceUi();
  }

}

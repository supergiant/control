import { Component, OnInit }                                          from '@angular/core';
import { ActivatedRoute }                                             from "@angular/router";
import { from, Observable, of }                                       from "rxjs";
import { filter, first, map, pluck, publishReplay, share, switchMap } from "rxjs/operators";
import { Supergiant }                                                 from "../../../shared/supergiant/supergiant.service";

@Component({
  selector: 'add-node',
  templateUrl: './add-node.component.html',
  styleUrls: ['./add-node.component.css']
})
export class AddNodeComponent implements OnInit {
  machines = [{
    machineType: null,
    role: "Master",
    qty: 1
  }];

  machineSizes$: Observable<any>;

  constructor(
    private supergiant: Supergiant,
    private route: ActivatedRoute,
  ) {
  }

  ngOnInit() {
    const clusterName = this.route.snapshot.params.id;

    const cluster$ = this.supergiant.Kubes.get(clusterName);

    const cloudAccountName$ = cluster$.pipe(
      pluck('accountName')
    );

    const firstNode$ = cluster$.pipe(
      pluck('nodes'),
      switchMap(nodes => from(Object.values(nodes))),
      first(),
      // publishReplay(),
    );

    const provider$ = firstNode$.pipe(
      pluck('provider')
    );


    this.machineSizes$ = cloudAccountName$.pipe(
      // filter(accountType => accountType === 'digitalocean' ),
      switchMap(accountName => this.supergiant.CloudAccounts.getRegions(accountName)),
      pluck('sizes'),
      map(sizes => Object.keys(sizes))
    );





    //
    const machineTypes$ = provider$.pipe(
      filter(region => region === 'digitalocean'),
      map(region => {
        return region;
      }))
      .subscribe(val => {
      console.log('subprovider', val);
    });
  }

}

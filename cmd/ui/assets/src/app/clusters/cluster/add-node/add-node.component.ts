import { HttpClient }                                                      from "@angular/common/http";
import { Component, OnDestroy, OnInit }                                    from '@angular/core';
import { ActivatedRoute }                                                  from "@angular/router";
import { combineLatest, from, Observable, of, Subject, Subscription }      from "rxjs";
import { distinctUntilChanged, filter, first, map, pluck, switchMap, tap } from "rxjs/operators";
import { Supergiant }                                                      from "../../../shared/supergiant/supergiant.service";
import { NodeProfileService }                                              from "../../node-profile.service";

@Component({
  selector: 'add-node',
  templateUrl: './add-node.component.html',
  styleUrls: ['./add-node.component.scss']
})
export class AddNodeComponent implements OnInit, OnDestroy {
  clusterName: string;

  machines = [{
    machineType: null,
    role: "Node",
    qty: 1
  }];

  machineSizes$: Observable<any>;
  provider$: Observable<any>;
  provider: string;
  providerSubj: Subject<any>;
  awsAvalableZones$: Observable<any>;

  clusterOptions = {
    archs: ["amd64"],
    flannelVersions: ["0.10.0"],
    operatingSystems: ["linux"],
    networkTypes: ["vxlan"],
    ubuntuVersions: ["xenial"],
    helmVersions: ["2.8.0"],
    dockerVersions: ["17.06.0"],
    K8sVersions: ["1.11.1"],
    rbacEnabled: [false]
  };

  subscriptions: Subscription;

  constructor(
    private supergiant: Supergiant,
    private route: ActivatedRoute,
    private http: HttpClient,
    private nodesService: NodeProfileService,
  ) {
    this.providerSubj = new Subject<string>();
    this.subscriptions = new Subscription();
  }

  ngOnInit() {
    this.clusterName = this.route.snapshot.params.id;

    const cluster$ = this.supergiant.Kubes.get(this.clusterName);

    const region$ = cluster$.pipe(
      distinctUntilChanged(),
      pluck('region'),
    );


    const cloudAccountName$ = cluster$.pipe(
      pluck('accountName')
    );

    const firstNode$ = cluster$.pipe(
      pluck('nodes'),
      switchMap(nodes => from(Object.values(nodes))),
      first()
    );

    this.provider$ = firstNode$.pipe(
      distinctUntilChanged(),
      pluck('provider'),
      tap(provider => {
        this.providerSubj.next(provider);
      })
    );

    this.subscriptions.add(
      this.provider$
        .subscribe(provider => this.provider = provider)
    );


    // Digital Ocean machine sizes
    const DOmachineSizes$ = cloudAccountName$.pipe(
      switchMap(accountName => this.supergiant.CloudAccounts.getRegions(accountName)),
      pluck('sizes'),
      map(sizes => Object.keys(sizes)),
    );

    const awsMachineSizes$ = combineLatest(this.provider$, region$).pipe(
      switchMap(([name, region]) =>
        combineLatest(
          of(name), of(region), this.supergiant.CloudAccounts.getAwsAvailabilityZones(name, region)
        )
      ),
      switchMap(([name, region, awsZones]) =>
        this.supergiant.CloudAccounts.getAwsMachineTypes(name, region, awsZones[0])
      ),
      pluck('sizes'),
      map(sizes => Object.keys(sizes)),
    );


    this.subscriptions.add(
      this.providerSubj.pipe(
        filter(provider => provider === 'digitalocean'),
        switchMap(() => DOmachineSizes$),
      ).subscribe(sizes => this.machineSizes$ = of(sizes))
    );

    this.subscriptions.add(
      this.providerSubj.pipe(
        filter(provider => provider === 'aws'),
        switchMap(_ => awsMachineSizes$),
      ).subscribe(sizes => this.machineSizes$ = of(sizes))
    );

  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  machineSizeChange(e) {
    console.log(e);
    console.log(this.machines);
  }

  addBlankMachine() {
    this.machines.push({
      machineType: null,
      role: 'Node',
      qty: 1
    })
  }

  finish() {
    const nodes = this.nodesService.compileProfiles(this.provider, this.machines, "Node");

    // TODO  move to service
    const url = `/v1/api/kubes/${this.clusterName}/nodes`;
    this.http.post(url, nodes)
      .subscribe(console.log);
  }
}

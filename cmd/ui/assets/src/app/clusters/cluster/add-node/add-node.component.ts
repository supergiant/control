import { HttpClient }                                                      from "@angular/common/http";
import {
  Component,
  OnDestroy,
  OnInit
}                                                                          from '@angular/core';
import { ActivatedRoute, Router }                                          from "@angular/router";
import { combineLatest, from, Observable, of, Subject, Subscription, zip } from "rxjs";
import {
  catchError,
  distinctUntilChanged,
  filter,
  first,
  map,
  mergeMap,
  pluck,
  switchMap,
  take,
  tap
}                                                                          from "rxjs/operators";
import { Notifications }                                                   from "../../../shared/notifications/notifications.service";
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
  isProcessing: boolean;

  constructor(
    private supergiant: Supergiant,
    private route: ActivatedRoute,
    private http: HttpClient,
    private nodesService: NodeProfileService,
    private notifications: Notifications,
    private router: Router,

  ) {
    this.providerSubj = new Subject<string>();
    this.subscriptions = new Subscription();
  }

  ngOnInit() {
    this.clusterName = this.route.snapshot.params.id;

    const cluster$ = this.supergiant.Kubes.get(this.clusterName);

    const region$ = cluster$.pipe(
      tap(p => {
        console.log('prov', p);
      }),
      pluck('region'),
      distinctUntilChanged()
    );


    const cloudAccountName$ = cluster$.pipe(
      pluck('accountName')
    );

    const firstNode$ = cluster$.pipe(
      pluck('masters'),
      switchMap(nodes => from(Object.values(nodes))),
      first()
    );

    this.provider$ = firstNode$.pipe(
      pluck('provider'),
      distinctUntilChanged(),
      tap(provider => {
        this.providerSubj.next(provider);
      })
    );

    // Digital Ocean machine sizes
    const DOmachineSizes$ = cloudAccountName$.pipe(
      switchMap(accountName => this.supergiant.CloudAccounts.getRegions(accountName)),
      pluck('sizes'),
      map(sizes => Object.keys(sizes)),
    );


    const awsMachineSizes$ = zip(this.provider$, region$).pipe(
      take(1),
      switchMap(([name, region]) =>
        combineLatest(
          of(name),
          of(region),
          this.supergiant.CloudAccounts.getAwsAvailabilityZones(name, region).pipe(
            // TODO: error handling
            catchError(e => of(e))
          )
        )
      ),
      mergeMap(([name, region, awsZones]) =>
        this.supergiant.CloudAccounts.getAwsMachineTypes(name, region, awsZones[0])
      ),
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
        first(),
        switchMap(_ => awsMachineSizes$)
      ).subscribe(sizes => this.machineSizes$ = of(sizes))
    );

    this.subscriptions.add(
      this.provider$
        .subscribe(provider => this.provider = provider)
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

  deleteMachine(idx) {
    if(this.machines.length === 1) return;

    this.machines.splice(idx, 1);
  }


  finish() {
    this.isProcessing = true;

    const nodes = this.nodesService.compileProfiles(this.provider, this.machines, "Node");
    // TODO  move to service
    const url = `/v1/api/kubes/${this.clusterName}/nodes`;

    this.http.post(url, nodes).pipe(
      catchError(error => {
        this.notifications.display('error', 'Error', error.statusText);
        return of(new ErrorEvent(error));
      })

    )
      .subscribe(result => {
        this.isProcessing = false;

        if (result instanceof ErrorEvent) return;

        this.notifications.display(
          'success',
          'Success!',
          'Your request is being processed'
        );
        this.router.navigate([`clusters/${this.clusterName}`]);
      });
  }
}

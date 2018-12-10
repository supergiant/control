import { HttpClient }                                                from "@angular/common/http";
import {
  Component,
  OnDestroy,
  OnInit
}                                                                    from '@angular/core';
import {
  ActivatedRoute,
  Router,
}                                                                    from "@angular/router";
import { combineLatest, Observable, of, Subject, Subscription, zip } from "rxjs";
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
}                                                                    from "rxjs/operators";
import { Notifications }                                             from "../../../shared/notifications/notifications.service";
import { Supergiant }                                                from "../../../shared/supergiant/supergiant.service";
import { CLUSTER_OPTIONS }                                           from "../../new-cluster/cluster-options.config";
import { NodeProfileService }                                        from "../../node-profile.service";

@Component({
  selector: 'add-node',
  templateUrl: './add-node.component.html',
  styleUrls: ['./add-node.component.scss']
})
export class AddNodeComponent implements OnInit, OnDestroy {
  clusterId: number;
  clusterName: number;

  machines = [{
    machineType: null,
    role: "Node",
    qty: 1,
    availabilityZone: '',
  }];

  machineSizes$: Observable<any>;
  provider$: Observable<any>;
  provider: string;
  providerSubj: Subject<any>;
  validMachinesConfig = false;
  displayMachineConfigError = false;

  clusterOptions = CLUSTER_OPTIONS;

  subscriptions: Subscription;
  isProcessing: boolean;
  availabilityZones: string[];
  selectedAZSubj: Subject<string>;
  isLoadingMachineTypes: boolean;
  machineTypesFilter: string = '';

  constructor(
    private supergiant: Supergiant,
    private route: ActivatedRoute,
    private http: HttpClient,
    private nodesService: NodeProfileService,
    private notifications: Notifications,
    private router: Router,
  ) {
    this.providerSubj = new Subject<string>();
    this.selectedAZSubj = new Subject<string>();
    this.subscriptions = new Subscription();
  }

  ngOnInit() {
    this.clusterId = this.route.snapshot.params.id;

    const cluster$ = this.supergiant.Kubes.get(this.clusterId);

    cluster$.subscribe(cluster => this.clusterName = cluster.name);

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

    this.provider$ = cluster$.pipe(
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


    const awsMachineSizes$: Observable<string[]> = zip(cloudAccountName$, region$).pipe(
      take(1),
      switchMap(([name, region]) =>
        combineLatest(
          of(name),
          of(region),
          this.supergiant.CloudAccounts.getAwsAvailabilityZones(name, region).pipe(
            tap(availabilityZones => this.availabilityZones = availabilityZones),
            // TODO: error handling
            catchError(e => of(e))
          )
        )
      ),
      // fetch machine types after az change
      mergeMap(([name, region]) => {
        return combineLatest(
          of({name, region}), this.selectedAZSubj,
          (params, awsZone) => [params.name, params.region, awsZone]
        );
      }),
      switchMap(([name, region, awsZone]) => this.supergiant.CloudAccounts.getAwsMachineTypes(name, region, awsZone))
    );

    const gceMachineSizes$: Observable<string[]> = zip(cloudAccountName$, region$).pipe(
      take(1),
      switchMap(([name, region]) =>
        combineLatest(
          of(name),
          of(region),
          this.supergiant.CloudAccounts.getGCEAvailabilityZones(name, region).pipe(
            tap(availabilityZones => this.availabilityZones = availabilityZones),
            // TODO: error handling
            catchError(e => of(e))
          )
        )
      ),
      // fetch machine types after az change
      mergeMap(([name, region]) => {
        return combineLatest(
          of({name, region}), this.selectedAZSubj,
          (params, gceZone) => [params.name, params.region, gceZone]
        );
      }),
      switchMap(([name, region, gceZone]) => this.supergiant.CloudAccounts.getGCEMachineTypes(name, region, gceZone))
    );

    this.subscriptions.add(
      this.providerSubj.pipe(
        filter(provider => provider === 'digitalocean'),
        switchMap(() => DOmachineSizes$),
      ).subscribe(sizes => this.machineSizes$ = of(sizes.sort()))
    );

    this.subscriptions.add(
      this.providerSubj.pipe(
        filter(provider => provider === 'aws'),
        first(),
        switchMap(_ => awsMachineSizes$)
      ).subscribe((sizes) => {
        this.isLoadingMachineTypes = false;
        this.machineSizes$ = of(sizes.sort());
      })
    );

    this.subscriptions.add(
      this.providerSubj.pipe(
        filter(provider => provider === 'gce'),
        first(),
        switchMap(_ => gceMachineSizes$)
      ).subscribe((sizes) => {
        this.isLoadingMachineTypes = false;
        this.machineSizes$ = of(sizes.sort());
      })
    );

    this.subscriptions.add(
      this.provider$
        .subscribe(provider => this.provider = provider)
    );

  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  addBlankMachine() {
    const lastMachine = this.machines[this.machines.length - 1];
    const newMachine = Object.assign({}, lastMachine);

    this.machines.push(newMachine);

    this.checkAndSetValidMachineConfig();
  }

  deleteMachine(idx) {
    if (this.machines.length === 1) return;

    this.machines.splice(idx, 1);

    this.checkAndSetValidMachineConfig();
  }

  validMachine(machine) {
    if (machine.machineType && machine.role && (typeof(machine.qty) == "number")) {
      return true;
    } else {
      return false;
    }
  }

  checkAndSetValidMachineConfig() {
    this.validMachinesConfig = this.machines.every(this.validMachine);

    if (this.validMachinesConfig) {
      this.displayMachineConfigError = false;
    } else {
      this.displayMachineConfigError = true;
    }
  }

  onAzChange(az) {
    this.isLoadingMachineTypes = true;
    this.selectedAZSubj.next(az);
  }

  finish() {

    if (this.validMachinesConfig) {
      this.isProcessing = true;
      this.displayMachineConfigError = false;

      const nodes = this.nodesService.compileProfiles(this.provider, this.machines, "Node");
      // TODO  move to service
      const url = `/v1/api/kubes/${this.clusterId}/nodes`;

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
          this.router.navigate([`clusters/${this.clusterId}`]);
        });

    } else {
      this.isProcessing = false;
      this.displayMachineConfigError = true;
    }
  }

  filterCallback = (val) => {
    if (this.machineTypesFilter === '') {
      return val;
    }

    return val.indexOf(this.machineTypesFilter) > -1;
  }

}

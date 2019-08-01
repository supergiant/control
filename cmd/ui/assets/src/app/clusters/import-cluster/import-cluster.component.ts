import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ValidatorFn, AbstractControl } from '@angular/forms';
import { Router } from '@angular/router';
import { Observable, of } from 'rxjs';
import { map, catchError } from 'rxjs/operators';

import { Notifications } from '../../shared/notifications/notifications.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { CLUSTER_OPTIONS } from '../new-cluster/cluster-options.config'

interface CloudAccount {
  name: string,
  provider: string,
  credentials: Object
};

@Component({
  selector: 'import-cluster',
  templateUrl: './import-cluster.component.html',
  styleUrls: ['./import-cluster.component.scss']
})
export class ImportClusterComponent implements OnInit {

  public cloudAccounts$: Observable<any>;
  public selectedCloudAccount: string;
  public availableRegions: Array<CloudAccount>;
  public importForm: FormGroup;
  public unavailableClusterNames = new Set;
  availableRegionNames: Array<string>;
  regionsFilter = '';
  regionsLoading = false;
  importing = false;
  clusterOptions = CLUSTER_OPTIONS;

  constructor(
    private supergiant: Supergiant,
    private notifications: Notifications,
    private router: Router,
    private formBuilder: FormBuilder,
  ) { }

  ngOnInit() {
    this.updateUsedClusterNames();
    this.cloudAccounts$ = this.supergiant.CloudAccounts.get().pipe(
      // importing only supported for aws currently
      map(cloudAccounts => cloudAccounts.sort().filter(c => c.provider == "aws")),
    );

    this.initImportForm();
  }

  intToHex(int: number) {
    return ("0" + int.toString(16)).substr(-2)
  }

  generateId(length?: number) {
    let arr = new Uint8Array((length || 40) / 2);
    window.crypto.getRandomValues(arr);
    return Array.from(arr, this.intToHex).join('')
  }

  updateUsedClusterNames() {
    this.supergiant.Kubes.get().subscribe(
      clusters => clusters.map(c => this.unavailableClusterNames.add(c.name)),
      err => console.error(err),
    );
  }

  initImportForm() {
    this.importForm = this.formBuilder.group({
      clusterName: ["", [
        Validators.required,
        this.uniqueClusterName(this.unavailableClusterNames),
        Validators.maxLength(12),
        Validators.pattern('^[A-Za-z]([-A-Za-z0-9\-]*[A-Za-z0-9\-])?$')
      ]],
      cloudAccount: ["", Validators.required],
      region: ["", Validators.required],
      arch: ["amd64", Validators.required],
      dockerVersion: ["18.06.3", Validators.required],
      K8SVersion: ["1.15.1", Validators.required],
      rbacEnabled: [true, Validators.required],
      kubeconfig: ["", Validators.required],
      publicKey: ["", Validators.required],
      privateKey: ["", Validators.required],
    })
  }

  importCluster(form: FormGroup) {
    if (form.valid) {
      const clusterId = this.generateId(8);
      const importClusterData = {
        clusterName: form.value.clusterName,
        cloudAccountName: form.value.cloudAccount.name,
        kubeconfig: form.value.kubeconfig,
        publicKey: form.value.publicKey,
        privateKey: form.value.privateKey,
        profile: {
          id: clusterId,
          region: form.value.region,
          arch: form.value.arch,
          operatingSystem: "linux",
          ubuntuVersion: "xenial",
          dockerVersion: form.value.dockerVersion,
          K8SVersion: form.value.K8SVersion,
          rbacEnabled: form.value.rbacEnabled,
        }
      };
      this.importing = true;
      this.supergiant.Kubes.import(importClusterData).subscribe(
        res => {
          this.success(importClusterData.clusterName);
          this.router.navigate(['/clusters/', res.clusterId]);
        },
        err => {
          this.importing = false;
          this.error(err);
          console.error(err);
        }
      )
    }
  }

  uniqueClusterName(unavailableNames: Set<string>): ValidatorFn {
    return (name: AbstractControl): { [key: string]: any } | null => {
      return unavailableNames.has(name.value) ? { 'nonUniqueName': { value: name.value } } : null;
    };
  }

  selectCloudAccountAndGetRegions(cloudAccount: CloudAccount) {
    this.selectedCloudAccount = cloudAccount.name;
    this.regionsLoading = true;
    const regions$ = this.supergiant.CloudAccounts.getRegions(cloudAccount.name)
      .pipe(
        catchError(err => {
          this.regionsLoading = false;
          console.error(err);
          this.error(err);
          return of("");
        }
      )
    )

    regions$.subscribe(res => {
      this.availableRegions = res.regions;
      this.availableRegionNames = this.availableRegions.map(n => n.name);
      this.regionsLoading = false;
    });
  }

  error(err) {
    this.notifications.display(
      "error",
      "Error",
      err.error.userMessage
    )
  }

  success(clusterName) {
    this.notifications.display(
      "success",
      "Success!",
      "Importing " + clusterName
    )
  }

  get clusterName() {
    return this.importForm.get("clusterName");
  }

}

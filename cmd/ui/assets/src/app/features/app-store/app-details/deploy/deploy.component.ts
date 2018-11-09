import { Component, Inject, OnInit }          from '@angular/core';
import { FormBuilder, FormGroup, Validators } from "@angular/forms";
import { State }                              from "../../../../reducers";
import { select, Store }                      from "@ngrx/store";
import { Chart, selectAppDetails }            from "../../../apps/apps.reducer";
import { Observable, of }                     from "rxjs";
import { HttpClient }                         from "@angular/common/http";
import { MAT_DIALOG_DATA, MatDialogRef }      from "@angular/material";
import { catchError }                         from "rxjs/operators";
import { Notifications }                      from "app/shared/notifications/notifications.service";
import { ActivatedRoute, Router }             from "@angular/router";


@Component({
  selector: 'deploy',
  templateUrl: './deploy.component.html',
  styleUrls: ['./deploy.component.scss']
})
export class DeployComponent implements OnInit {

  deployForm: FormGroup;
  currentChart$: Observable<Chart>;
  clusters$: Observable<any>;
  isProcessing: boolean;

  constructor(
    private formBuilder: FormBuilder,
    private store: Store<State>,
    private http: HttpClient,
    @Inject(MAT_DIALOG_DATA) public data: any,
    private notifications: Notifications,
    private route: ActivatedRoute,
    public router: Router,
    private dialogRef: MatDialogRef<DeployComponent>,
  ) {
  }

  ngOnInit() {
    this.deployForm = this.formBuilder.group({
      clusterName: ['', { disabled: true }, Validators.required],
      name: [''],
      namespace: ['default'],
      chartName: [''],
      chartVersion: [''],
      repoName: [''],
      values: [''],
    });

    this.currentChart$ = this.store.pipe(select(selectAppDetails));
    this.clusters$ = this.http.get('/v1/api/kubes');
    this.setDefaultFormValues();
  }

  submitForm() {
    const formValue = this.deployForm.getRawValue();
    this.deployForm.disable();
    this.isProcessing = true;


    this.http.post(`/v1/api/kubes/${formValue.clusterName}/releases`, formValue).pipe(
      catchError(error => {
        this.notifications.display('error', 'Error', error.statusText);
        return of(new ErrorEvent(error));
      })
    ).subscribe(result => {
      this.isProcessing = false;
      this.deployForm.enable();
      this.disableUnusedFields();

      if (result instanceof ErrorEvent) {
        return;
      }
      this.router.navigate(['apps']);
      this.notifications.display('success', 'Success', 'App is being deployed!');

      this.dialogRef.close()
    });
  }

  private setDefaultFormValues() {
    this.currentChart$.subscribe(currentChart => {
      const repoName = this.data.routeParams.repo;
      const chartName = currentChart.metadata.name;

      this.deployForm.patchValue({
        chartName, repoName, values: currentChart.values
      });
    });

    this.disableUnusedFields();
  }

  private disableUnusedFields() {
    this.deployForm.controls.chartName.disable();
    this.deployForm.controls.chartVersion.disable();
  }
}

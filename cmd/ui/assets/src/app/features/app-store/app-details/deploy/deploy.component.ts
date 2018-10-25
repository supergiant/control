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
import { Router }                             from "@angular/router";


@Component({
  selector: 'deploy',
  templateUrl: './deploy.component.html',
  styleUrls: [ './deploy.component.scss' ]
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
    public router: Router,
    private dialogRef: MatDialogRef<DeployComponent>,
  ) {
  }

  ngOnInit() {
    this.deployForm = this.formBuilder.group({
      clusterName: [ '', { disabled: true }, Validators.required ],
      name: [ '' ],
      namespace: [ '' ],
      chartName: [ '' ],
      chartVersion: [ '' ],
      repoName: [ '' ],
      values: [ '' ]
    });

    this.currentChart$ = this.store.pipe(select(selectAppDetails));
    this.clusters$ = this.http.get('/v1/api/kubes');
    this.setDefaultFormValues();
  }

  submitForm() {
    const formValue = this.deployForm.getRawValue();
    this.deployForm.disable();
    this.isProcessing = true;

    this.http.post(`/v1/api/${formValue.clusterName}/releases`, formValue).pipe(
      catchError(error => {
        this.notifications.display('error', 'Error', error.statusText);
        return of(new ErrorEvent(error));
      })
    ).subscribe(result => {
      this.isProcessing = false;
      this.deployForm.enable();

      if (result instanceof ErrorEvent) {
        return;
      }
      this.router.navigate([ 'apps' ]);
      this.notifications.display('success', 'Success', 'App is being deployed!');

      this.dialogRef.close()
    });
  }

  private setDefaultFormValues() {
    this.currentChart$.subscribe(currentChart => {
      this.deployForm.patchValue({
        chartName: currentChart.name,
        repoName: currentChart.repo,
      });
    })
  }

}

import { Component, Inject, OnInit }          from '@angular/core';
import { FormBuilder, FormGroup, Validators } from "@angular/forms";
import { State }                              from "../../../../reducers";
import { select, Store }                      from "@ngrx/store";
import { Chart, selectAppDetails }            from "../../../apps/apps.reducer";
import { Observable }                         from "rxjs";
import { HttpClient }                         from "@angular/common/http";
import { MAT_DIALOG_DATA }                    from "@angular/material";

@Component({
  selector: 'deploy',
  templateUrl: './deploy.component.html',
  styleUrls: [ './deploy.component.scss' ]
})
export class DeployComponent implements OnInit {

  deployForm: FormGroup;
  currentChart$: Observable<Chart>;
  clusters$: Observable<any>;
  selecteCluster: string;

  constructor(
    private formBuilder: FormBuilder,
    private store: Store<State>,
    private http: HttpClient,
    @Inject(MAT_DIALOG_DATA) public data: any,
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

    this.http.post(`/v1/api/${formValue.clusterName}/releases`, formValue)
        .subscribe(result => {
          console.log(result);
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

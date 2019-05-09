import { Component, OnInit }      from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { HttpClient }             from '@angular/common/http';
import { Router } from '@angular/router';
import { catchError }             from 'rxjs/operators';
import { of }                     from 'rxjs';
import { MatDialogRef } from '@angular/material';
import { Notifications }          from '../../../shared/notifications/notifications.service';

@Component({
  selector: 'app-apps-add',
  templateUrl: './apps-add.component.html',
  styleUrls: [ './apps-add.component.scss' ]
})
export class AppsAddComponent implements OnInit {
  addRepositoryForm: FormGroup;
  isProcessing: boolean = false;
  disableSubmit: boolean ;

  constructor(
    private formBuilder: FormBuilder,
    private http: HttpClient,
    private notifications: Notifications,
    private router: Router,
    private dialogRef: MatDialogRef<AppsAddComponent>
  ) {
  }

  ngOnInit() {
    const pattern = new RegExp('^[a-zA-Z0-9-_]+$');
    this.addRepositoryForm = this.formBuilder.group({
      name: [ '', [Validators.required, Validators.pattern(pattern)] ],
      url: [ '', Validators.required ],
    });
    this.updateValidity();
  }

  get name() {
    return this.addRepositoryForm.get("name");
  }

  updateValidity() {
    this.addRepositoryForm.updateValueAndValidity();
    this.disableSubmit = (this.isProcessing || this.addRepositoryForm.invalid);
  }

  addRepository() {
    this.isProcessing = true;
    this.addRepositoryForm.disable();

    this.http.post(
      '/v1/api/helm/repositories',
      this.addRepositoryForm.getRawValue()
    ).pipe(
      catchError(error => {
        this.notifications.display('error', '', error.statusText);
        return of(new ErrorEvent(error));
      })
    ).subscribe(result => {
      this.isProcessing = false;
      const repoName = this.addRepositoryForm.value.name;
      this.addRepositoryForm.enable();
      if (!(result instanceof ErrorEvent)) {
        this.router.navigate(['/catalog/', repoName]);
        this.dialogRef.close()
      }
    });
  }
}

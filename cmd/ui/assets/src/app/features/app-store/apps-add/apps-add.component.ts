import { Component, OnInit }      from '@angular/core';
import { FormBuilder, FormGroup } from "@angular/forms";
import { HttpClient }             from "@angular/common/http";
import { catchError }             from "rxjs/operators";
import { of }                     from "rxjs";
import { Notifications }          from "../../../shared/notifications/notifications.service";

@Component({
  selector: 'app-apps-add',
  templateUrl: './apps-add.component.html',
  styleUrls: [ './apps-add.component.scss' ]
})
export class AppsAddComponent implements OnInit {
  addRepositoryForm: FormGroup;
  isProcessing: boolean;

  constructor(
    private formBuilder: FormBuilder,
    private http: HttpClient,
    private notifications: Notifications,
  ) {
  }

  ngOnInit() {
    this.addRepositoryForm = this.formBuilder.group({
      name: [ '' ],
      url: [ '' ],
    });

    this.http.get('/v1/api/helm/repositories').subscribe( res => {
      console.log(res);
    });
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
        console.error(error);
        return of(error);
      })
    ).subscribe(res => {
      this.isProcessing = false;
      this.addRepositoryForm.enable();
      // TODO
      window.location.reload();
    });
  }
}

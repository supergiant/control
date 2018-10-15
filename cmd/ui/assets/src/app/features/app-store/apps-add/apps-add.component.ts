import { Component, OnInit }      from '@angular/core';
import { FormBuilder, FormGroup } from "@angular/forms";
import { HttpClient }             from "@angular/common/http";

@Component({
  selector: 'app-apps-add',
  templateUrl: './apps-add.component.html',
  styleUrls: [ './apps-add.component.scss' ]
})
export class AppsAddComponent implements OnInit {
  addRepositoryForm: FormGroup;

  constructor(
    private formBuilder: FormBuilder,
    private http: HttpClient,
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
    this.http.post(
      '/v1/api/helm/repositories',
      this.addRepositoryForm.getRawValue()
    ).subscribe(res => {
      console.log(res);
    })
  }
}

import { Component, OnInit }      from '@angular/core';
import { FormBuilder, FormGroup } from "@angular/forms";

@Component({
  selector: 'app-apps-add',
  templateUrl: './apps-add.component.html',
  styleUrls: [ './apps-add.component.scss' ]
})
export class AppsAddComponent implements OnInit {
  addRepository: FormGroup;

  constructor(
    private formBuilder: FormBuilder,
  ) {
  }

  ngOnInit() {
    this.addRepository = this.formBuilder.group({
      name: [ '' ],
      url: [ '' ],
    })
  }

}

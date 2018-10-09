import { Component, Input, OnInit } from '@angular/core';
import { Chart }                    from "../../apps/apps.reducer";

@Component({
  selector: 'apps-list',
  templateUrl: './apps-list.component.html',
  styleUrls: ['./apps-list.component.scss']
})
export class AppsListComponent implements OnInit {

  @Input() charts: Chart[];

  constructor() { }

  ngOnInit() {
  }

}

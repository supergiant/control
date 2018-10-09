import { Component, OnInit } from '@angular/core';
import { ActivatedRoute }    from "@angular/router";

@Component({
  selector: 'app-details',
  templateUrl: './app-details.component.html',
  styleUrls: [ './app-details.component.scss' ]
})
export class AppDetailsComponent implements OnInit {

  constructor(private route: ActivatedRoute,) {}

  ngOnInit() {
    let id = this.route.snapshot.paramMap.get('repo');
    let chart = this.route.snapshot.paramMap.get('chart');
  }

}

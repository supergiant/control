import { Component, OnInit } from '@angular/core';
import { ActivatedRoute }    from '@angular/router';

@Component({
  selector: 'breadcrumbs',
  templateUrl: './breadcrumbs.component.html',
  styleUrls: ['./breadcrumbs.component.scss']
})
export class BreadcrumbsComponent implements OnInit {
  repoName;
  chartName;

  constructor(
    private route: ActivatedRoute,
  ) {
  }

  ngOnInit() {
    const { repo, chart } = this.route.children[0].snapshot.params;
    this.repoName = repo;
    this.chartName = chart;
  }
}

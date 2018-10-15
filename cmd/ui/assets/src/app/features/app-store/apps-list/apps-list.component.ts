import { Component, Input, OnInit } from '@angular/core';
import { Chart }                    from "../../apps/apps.reducer";
import { ActivatedRoute, Router }   from "@angular/router";
import { HttpClient }               from "@angular/common/http";
import { Observable }               from "rxjs";

@Component({
  selector: 'apps-list',
  templateUrl: './apps-list.component.html',
  styleUrls: [ './apps-list.component.scss' ]
})
export class AppsListComponent implements OnInit {

  @Input() charts: Chart[] | Observable<any>;
  repo: string;

  constructor(
    private route: ActivatedRoute,
    private http: HttpClient,
    public router: Router,
  ) {
    this.repo = this.route.snapshot.paramMap.get('repo');
    this.charts = this.http.get(`/v1/api/helm/repositories/${this.repo}/charts`);
  }

  ngOnInit() {
    this.router.events.subscribe(() => {
      this.repo = this.route.snapshot.paramMap.get('repo');

      this.charts = this.http.get(`/v1/api/helm/repositories/${this.repo}/charts`);
    });

  }

}

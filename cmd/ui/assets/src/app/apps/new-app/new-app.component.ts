import { Component, OnInit, ViewEncapsulation, OnDestroy } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { ActivatedRoute, Router } from '@angular/router';
import { AppsModel } from '../apps.model';
import * as GenerateSchema from 'generate-schema';
import { Location } from '@angular/common';

@Component({
  selector: 'app-new-app',
  templateUrl: './new-app.component.html',
  styleUrls: ['./new-app.component.scss']
})
export class NewAppComponent implements OnInit, OnDestroy {
  private subscriptions = new Subscription();
  private appsModel = new AppsModel;
  private clusters = [];
  private chart: any;
  public model: any;
  public schema: any;
  private id: number;
  public loaded = false;

  constructor(
    private supergiant: Supergiant,
    private notifications: Notifications,
    private router: Router,
    private route: ActivatedRoute,
    private location: Location,
  ) { }

  ngOnInit() {
    this.id = this.route.snapshot.params.id;
    // quick hack, tucking get into get clusters so enum isnt empty
    this.getClusters();

  }

  get(id) {
    this.subscriptions.add(this.supergiant.HelmCharts.get(id).subscribe(
      (chart) => {
        if (chart.default_config) {
          // this is our model: the vars file provided by the chart.
          this.appsModel.app.model.config = JSON.parse(JSON.stringify(chart.default_config).replace(/\[\]/g, '["Enter-Info"]', ));
          // this.appsModel.app.model.config = chart.default_config;

        }

        this.appsModel.app.model.chart_name = chart.name;
        this.appsModel.app.model.chart_version = chart.version;
        // this needs to be called cluster in display
        // this.appsModel.app.model.kube_name = 'FIXME';
        this.appsModel.app.model.repo_name = chart.repo_name;
        // We dynamically generate a schema from the vars file.
        // TODO: Note - we need to add a sniffer here to look for a schema file
        // and any icon images in the chart. If they exist, we should use them instead of the generated one.


        this.appsModel.app.schema.properties.config = GenerateSchema.json(this.appsModel.app.model.config);

        if (this.clusters.length) {
          this.appsModel.app.schema.properties.kube_name.enum = this.clusters;
        } else {
          this.appsModel.app.schema.properties.kube_name.enum = ['No Kubes found'];
        }



        this.model = this.appsModel.app.model;
        this.schema = this.appsModel.app.schema;

        this.loaded = true;
      },
      (err) => { console.log(err); }));
  }

  create(model) {
    this.subscriptions.add(this.supergiant.HelmReleases.create(model).subscribe(
      (data) => {
        this.success(model);
        this.router.navigate(['/apps']);
      },
      (err) => { this.error(model, err); }));
  }

  success(model) {
    this.notifications.display(
      'success',
      'App: ' + model.chart_name,
      'Deployed...',
    );
  }

  error(model, data) {
    this.notifications.display(
      'error',
      'App: ' + model.chart_name,
      'Error:' + data);
  }

  getClusters() {
    this.subscriptions.add(this.supergiant.Kubes.get().subscribe(
      (kubes) => {
        const clusters: Array<any> = [];

        for (const kube of kubes.items) {
          clusters.push(kube.name);
        }

        this.clusters = clusters;

        this.get(this.id);
      },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  goBack() {
    this.location.back();
  }
}

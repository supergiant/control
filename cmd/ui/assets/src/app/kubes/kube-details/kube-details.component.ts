import { Component, OnInit, OnDestroy, ViewChild, ElementRef, ViewContainerRef } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { ActivatedRoute, Router } from '@angular/router';
import { Observable } from 'rxjs/Observable';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { ChartsModule, BaseChartDirective } from 'ng2-charts';
import { SystemModalService } from '../../shared/system-modal/system-modal.service';
import { LoginComponent } from '../../login/login.component';
import { DomSanitizer, SafeResourceUrl } from '@angular/platform-browser';

@Component({
  selector: 'app-kube-details',
  templateUrl: './kube-details.component.html',
  styleUrls: ['./kube-details.component.css']
})
export class KubeDetailsComponent implements OnInit, OnDestroy {
  id: number;
  subscriptions = new Subscription();
  kube: any;
  url: string;
  private tabSet: ViewContainerRef;
  @ViewChild('iframe') iframe: ElementRef;
  @ViewChild('t') ngbTabSet;
  public isLoading: Boolean;
  public secureSrc: SafeResourceUrl;
  public planets = [];
  public planetName: string;
  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private supergiant: Supergiant,
    private notifications: Notifications,
    private chartsModule: ChartsModule,
    private systemModalService: SystemModalService,
    public loginComponent: LoginComponent,
    private sanitizer: DomSanitizer,
  ) { }

  // CPU Usage
  public cpuChartData: Array<any> = [];
  public cpuChartOptions: any = {
    responsive: true
  };
  public cpuChartLabels: Array<any> = [];
  public cpuChartType = 'line';

  // RAM Usage
  public ramChartData: Array<any> = [];
  public ramChartOptions: any = {
    responsive: true
  };
  public ramChartLabels: Array<any> = [];
  public ramChartType = 'line';


  isDataAvailable = false;
  ngOnInit() {
    this.id = this.route.snapshot.params.id;
    this.getKube();
  }

  openSystemModal(message) {
    this.systemModalService.openSystemModal(message);
  }

  getKube() {
    this.subscriptions.add(Observable.timer(0, 10000)
      .switchMap(() => this.supergiant.Kubes.get(this.id)).subscribe(
      (kube) => {
        this.kube = kube;
        if (this.kube.extra_data && this.kube.extra_data.cpu_usage_rate && this.kube.extra_data.kube_cpu_capacity) {
          this.isDataAvailable = true;
          this.cpuChartData = [
            { label: 'CPU Usage', data: this.kube.extra_data.cpu_usage_rate.map((data) => data.value) },
            { label: 'CPU Capacity', data: this.kube.extra_data.kube_cpu_capacity.map((data) => data.value) },
            // this should be set to the length of largest array.
          ];
          this.ramChartLabels = this.kube.extra_data.cpu_usage_rate.map((data) => data.timestamp);

          this.ramChartData = [
            { label: 'RAM Usage', data: this.kube.extra_data.memory_usage.map((data) => data.value / 1073741824) },
            {
              label: 'RAM Capacity',
              data: this.kube.extra_data.kube_memory_capacity.map((data) => data.value / 1073741824)
            },
            // this should be set to the length of largest array.
          ];
          this.cpuChartLabels = this.kube.extra_data.memory_usage.map((data) => data.timestamp);
        }
      },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));

    // Get any planets
    this.subscriptions.add(Observable.timer(0, 20000)
      .switchMap(() => this.supergiant.KubeResources.get()).subscribe(
      (services) => {
        this.planets = services.items.filter(
          planet => {
            if (planet.resource.metadata.labels) {
              return planet.resource.metadata.labels['kubernetes.io/cluster-service'] === 'true';
            }
          }
        );
      },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  padArrayWithDefault(arr: any, n: number) {
    let tmpArr = [];
    tmpArr = arr.slice(0);
    while (tmpArr.length < n) {
      let count = 0;
      arr = tmpArr.slice(0);
      arr.reduce((previous, current, index) => {
        if (previous && tmpArr.length < n) {
          const average = (current + previous) / 2;
          tmpArr.splice(index + count, 0, average);
          count += 1;
        }
        return current;
      });
    }
    return tmpArr;
  }

  resetTabs(tab) {
    if (tab.nextId !== 'planetTab') {
      this.planetName = '';
    }
  }

  getIframeURL(name) {
    this.planetName = name;
    const service = '/api/v1/proxy/namespaces/kube-system/services/' + name;
    const basicAuth = this.kube.username + ':' + this.kube.password; // Can we send this somehow??
    this.url = 'https://' + this.kube.master_public_ip + service;
    this.secureSrc = this.sanitizer.bypassSecurityTrustResourceUrl(this.url);
    this.ngbTabSet.select('planetTab');
  }

  onIframeLoad() {
    const basicAuth = 'Basic ' + btoa(this.kube.username + ':' + this.kube.password);
    if (typeof this.iframe !== 'undefined') {
      this.iframe
        .nativeElement
        .contentWindow
        .postMessage('Authorization:', basicAuth);

      this.isLoading = false;
    }
  }

  goBack() {
    this.router.navigate(['/kubes']);
  }
  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

}

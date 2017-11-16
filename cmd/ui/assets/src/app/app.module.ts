// Modules
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgModule } from '@angular/core';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { HttpModule } from '@angular/http';
import { SimpleNotificationsModule } from 'angular2-notifications';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { SchemaFormModule, WidgetRegistry, DefaultWidgetRegistry } from 'angular2-schema-form';
import { JsonSchemaFormModule } from 'angular2-json-schema-form';
import { AppRoutingModule } from './app-routing.module';
import { ChartsModule } from 'ng2-charts';
import { NgxPaginationModule } from 'ngx-pagination';

// Components
import { AppComponent } from './app.component';
import { NavigationComponent } from './navigation/navigation.component';
import { UsersComponent } from './users/users.component';
import { KubesComponent } from './kubes/kubes.component';
import { NotificationsComponent } from './shared/notifications/notifications.component';
import { KubeComponent } from './kubes/kube/kube.component';
import { KubesHeaderComponent } from './kubes/kubes-header/kubes-header.component';
import { SessionsComponent } from './sessions/sessions.component';
import { CloudAccountsComponent } from './cloud-accounts/cloud-accounts.component';
import { LoadBalancersComponent } from './load-balancers/load-balancers.component';
import { NodesComponent } from './nodes/nodes.component';
import { PodsComponent } from './pods/pods.component';
import { ServicesComponent } from './services/services.component';
import { SessionsHeaderComponent } from './sessions/sessions-header/sessions-header.component';
import { ServicesHeaderComponent } from './services/services-header/services-header.component';
import { PodsHeaderComponent } from './pods/pods-header/pods-header.component';
import { NodesHeaderComponent } from './nodes/nodes-header/nodes-header.component';
import { LoadBalancersHeaderComponent } from './load-balancers/load-balancers-header/load-balancers-header.component';
import { CloudAccountsHeaderComponent } from './cloud-accounts/cloud-accounts-header/cloud-accounts-header.component';
import { CloudAccountComponent } from './cloud-accounts/cloud-account/cloud-account.component';
import { LoadBalancerComponent } from './load-balancers/load-balancer/load-balancer.component';
import { NodeComponent } from './nodes/node/node.component';
import { PodComponent } from './pods/pods/pod.component';
import { VolumesComponent } from './volumes/volumes.component';
import { VolumeComponent } from './volumes/volume/volume.component';
import { VolumesHeaderComponent } from './volumes/volumes-header/volumes-header.component';
import { ServiceComponent } from './services/service/service.component';
import { SessionComponent } from './sessions/session/session.component';
import { UserComponent } from './users/user/user.component';
import { UsersHeaderComponent } from './users/users-header/users-header.component';
import { SystemModalComponent } from './shared/system-modal/system-modal.component';
import { DropdownModalComponent } from './shared/dropdown-modal/dropdown-modal.component';
import { EditModalComponent } from './shared/edit-modal/edit-modal.component';
import { LoginComponent } from './login/login.component';
import { CookiesComponent } from './shared/cookies/cookies.component';
import { AppsComponent } from './apps/apps.component';
import { AppsHeaderComponent } from './apps/apps-header/apps-header.component';
import { Search } from './shared/search-pipe/search-pipe';
import { HelmAppComponent } from './apps/app/helm-app.component';
import { DeploymentComponent } from './apps/deployment/deployment.component';
import { RepoModalComponent } from './apps/repo-modal/repo-modal.component';
import { KubeDetailsComponent } from './kubes/kube-details/kube-details.component';
import { KubesListComponent } from './kubes/kubes-list/kubes-list.component';
import { NodeDetailsComponent } from './nodes/node-details/node-details.component';
import { NodesListComponent } from './nodes/nodes-list/nodes-list.component';
import { PodDetailsComponent } from './pods/pod-details/pod-details.component';
import { PodsListComponent } from './pods/pods-list/pods-list.component';
import { SessionDetailsComponent } from './sessions/session-details/session-details.component';
import { SessionsListComponent } from './sessions/sessions-list/sessions-list.component';
import { UserDetailsComponent } from './users/user-details/user-details.component';
import { UsersListComponent } from './users/users-list/users-list.component';
import { CloudAccountDetailsComponent } from './cloud-accounts/cloud-account-details/cloud-account-details.component';
import { CloudAccountsListComponent } from './cloud-accounts/cloud-accounts-list/cloud-accounts-list.component';
import { LoadBalancerDetailsComponent } from './load-balancers/load-balancer-details/load-balancer-details.component';
import { LoadBalancersListComponent } from './load-balancers/load-balancers-list/load-balancers-list.component';
import { ServiceDetailsComponent } from './services/service-details/service-details.component';
import { ServicesListComponent } from './services/services-list/services-list.component';
import { SupergiantComponent } from './shared/supergiant/supergiant.component';
import { AppsListComponent } from './apps/apps-list/apps-list.component';
import { AppDetailsComponent } from './apps/app-details/app-details.component';
import { DeploymentDetailsComponent } from './apps/deployment-details/deployment-details.component';
// Component Services
import { SessionsService } from './sessions/sessions.service';
import { CloudAccountsService } from './cloud-accounts/cloud-accounts.service';
import { KubesService } from './kubes/kubes.service';
import { UsersService } from './users/users.service';
import { NodesService } from './nodes/nodes.service';
import { PodsService } from './pods/pods.service';
import { AppsService } from './apps/apps.service';
import { VolumesService } from './volumes/volumes.service';
import { ServicesService } from './services/services.service';
import { LoadBalancersService } from './load-balancers/load-balancers.service';
import { Notifications } from './shared/notifications/notifications.service';
import { SystemModalService } from './shared/system-modal/system-modal.service';
import { DropdownModalService } from './shared/dropdown-modal/dropdown-modal.service';
import { EditModalService } from './shared/edit-modal/edit-modal.service';
import { CookieMonster } from './shared/cookies/cookies.service';
import { RepoModalService } from './apps/repo-modal/repo-modal.service';

// Supergiant API Services
import { Supergiant } from './shared/supergiant/supergiant.service';
import { UtilService } from './shared/supergiant/util/util.service';
import { Sessions } from './shared/supergiant/sessions/sessions.service';
import { Users } from './shared/supergiant/users/users.service';
import { CloudAccounts } from './shared/supergiant/cloud-accounts/cloud-accounts.service';
import { Kubes } from './shared/supergiant/kubes/kubes.service';
import { KubeResources } from './shared/supergiant/kube-resources/kube-resources.service';
import { Nodes } from './shared/supergiant/nodes/nodes.service';
import { LoadBalancers } from './shared/supergiant/load-balancers/load-balancers.service';
import { HelmRepos } from './shared/supergiant/helm-repos/helm-repos.service';
import { HelmCharts } from './shared/supergiant/helm-charts/helm-charts.service';
import { HelmReleases } from './shared/supergiant/helm-releases/helm-releases.service';
import { Logs } from './shared/supergiant/logs/logs.service';
import { AuthenticatedHttpService } from './shared/auth/authenticated-http-service.service';
import { Http } from '@angular/http';

import 'rxjs/add/operator/map';
import 'rxjs/add/operator/switchMap';
import 'rxjs/add/observable/of';
import 'rxjs/add/observable/timer';
import 'rxjs/add/observable/interval';
import 'rxjs/add/operator/take';

// ui2000
import { SystemComponent } from './system/system.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { ClustersComponent } from './clusters/clusters.component';
import { NewCloudAccountComponent } from './system/cloud-accounts/new-cloud-account/new-cloud-account.component';
// temporary 2000 name hack because of conflict
import { CloudAccount2000Component } from './system/cloud-accounts/cloud-account/cloud-account.component';
import { CloudAccounts2000Component } from './system/cloud-accounts/cloud-accounts.component';
import { Users2000Component } from './system/users/users.component';
import { EditCloudAccountComponent } from './system/cloud-accounts/edit-cloud-account/edit-cloud-account.component';
import { MainComponent } from './system/main/main.component';
import { HelmReposComponent } from './system/main/helm-repos/helm-repos.component';
import { NewClusterComponent } from './clusters/new-cluster/new-cluster.component';
import { ClusterComponent } from './clusters/cluster/cluster.component';







@NgModule({
  declarations: [
    AppComponent,
    VolumesComponent,
    VolumeComponent,
    VolumesHeaderComponent,
    NavigationComponent,
    UsersComponent,
    KubesComponent,
    KubeComponent,
    KubesHeaderComponent,
    SessionsComponent,
    CloudAccountsComponent,
    LoadBalancersComponent,
    NodesComponent,
    PodsComponent,
    ServicesComponent,
    SessionsHeaderComponent,
    ServicesHeaderComponent,
    PodsHeaderComponent,
    NodesHeaderComponent,
    LoadBalancersHeaderComponent,
    CloudAccountsHeaderComponent,
    CloudAccountComponent,
    LoadBalancerComponent,
    NodeComponent,
    PodComponent,
    ServiceComponent,
    SessionComponent,
    UserComponent,
    UsersHeaderComponent,
    NotificationsComponent,
    SystemModalComponent,
    DropdownModalComponent,
    EditModalComponent,
    LoginComponent,
    CookiesComponent,
    AppsComponent,
    AppsHeaderComponent,
    HelmAppComponent,
    DeploymentComponent,
    RepoModalComponent,
    KubeDetailsComponent,
    KubesListComponent,
    NodeDetailsComponent,
    NodesListComponent,
    PodDetailsComponent,
    PodsListComponent,
    SessionDetailsComponent,
    SessionsListComponent,
    UserDetailsComponent,
    UsersListComponent,
    CloudAccountDetailsComponent,
    CloudAccountsListComponent,
    LoadBalancerDetailsComponent,
    LoadBalancersListComponent,
    ServiceDetailsComponent,
    ServicesListComponent,
    Search,
    SupergiantComponent,
    AppDetailsComponent,
    AppsListComponent,
    DeploymentDetailsComponent,
    DashboardComponent,
    ClustersComponent,
    NewCloudAccountComponent,
    EditCloudAccountComponent,
    MainComponent,
    HelmReposComponent,
    NewClusterComponent,
    ClusterComponent,
    CloudAccount2000Component,
    CloudAccounts2000Component,
    SystemComponent,
    HelmReposComponent,
    Users2000Component,
  ],
  imports: [
    BrowserModule,
    NgbModule.forRoot(),
    AppRoutingModule,
    HttpModule,
    FormsModule,
    BrowserModule,
    BrowserAnimationsModule,
    SimpleNotificationsModule.forRoot(),
    ReactiveFormsModule,
    SchemaFormModule,
    ChartsModule,
    NgxPaginationModule,
    JsonSchemaFormModule,
  ],
  providers: [
    // Component Services
    KubesService,
    CloudAccountsService,
    SessionsService,
    UsersService,
    KubesService,
    NodesService,
    LoadBalancersService,
    PodsService,
    ServicesService,
    VolumesService,
    AppsService,
    RepoModalService,
    // Supergiant API Services
    Supergiant,
    UtilService,
    CloudAccounts,
    Sessions,
    Users,
    Kubes,
    KubeResources,
    Nodes,
    LoadBalancers,
    HelmRepos,
    HelmCharts,
    HelmReleases,
    Logs,
    // Other Shared Services
    { provide: WidgetRegistry, useClass: DefaultWidgetRegistry },
    Notifications,
    SystemModalService,
    DropdownModalService,
    EditModalService,
    CookieMonster,
    LoginComponent,
    { provide: Http, useClass: AuthenticatedHttpService },
  ],
  bootstrap: [AppComponent]
})

export class AppModule { }

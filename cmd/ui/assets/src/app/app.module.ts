// Modules
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgModule } from '@angular/core';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { HttpClientModule } from '@angular/common/http';
import { SimpleNotificationsModule } from 'angular2-notifications/dist';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import {
  MatAutocompleteModule,
  MatButtonModule,
  MatButtonToggleModule,
  MatCardModule,
  MatCheckboxModule,
  MatChipsModule,
  MatDatepickerModule,
  MatDialogModule,
  MatDividerModule,
  MatExpansionModule,
  MatGridListModule,
  MatIconModule,
  MatInputModule,
  MatListModule,
  MatMenuModule,
  MatNativeDateModule,
  MatPaginatorModule,
  MatProgressBarModule,
  MatProgressSpinnerModule,
  MatRadioModule,
  MatRippleModule,
  MatSelectModule,
  MatSidenavModule,
  MatSliderModule,
  MatSlideToggleModule,
  MatSnackBarModule,
  MatSortModule,
  MatStepperModule,
  MatTableModule,
  MatTabsModule,
  MatToolbarModule,
  MatTooltipModule } from '@angular/material';
import { SchemaFormModule, WidgetRegistry, DefaultWidgetRegistry } from 'ngx-schema-form';
import { MaterialDesignFrameworkModule } from 'angular6-json-schema-form';
import { AppRoutingModule } from './app-routing.module';
import { AceEditorModule } from 'ng2-ace-editor';
import { ChartsModule } from 'ng2-charts';
import { NgxPaginationModule } from 'ngx-pagination';
import { TitleCasePipe } from '@angular/common';
import { CommonModule } from '@angular/common';
import { ContextMenuModule } from 'ngx-contextmenu';


// Components
import { AppComponent } from './app.component';
import { NavigationComponent } from './navigation/navigation.component';
import { UsersComponent } from './users/users.component';
import { NotificationsComponent } from './shared/notifications/notifications.component';
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
import { Search } from './shared/search-pipe/search-pipe';
import { HelmAppComponent } from './apps/app/helm-app.component';
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
// Component Services
import { SessionsService } from './sessions/sessions.service';
import { CloudAccountsService } from './cloud-accounts/cloud-accounts.service';
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
import { AuthService } from './shared/supergiant/auth/auth.service';
import { AuthGuardService } from './shared/supergiant/auth/auth-guard.service';
import { TokenInterceptor } from './shared/supergiant/auth/token.interceptor';
import { Http } from '@angular/http';
import { HTTP_INTERCEPTORS } from '@angular/common/http';








// ui2000
import { SystemComponent } from './system/system.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { ClustersComponent } from './clusters/clusters.component';
import { NewCloudAccountComponent } from './system/cloud-accounts/new-cloud-account/new-cloud-account.component';
import { NgxDatatableModule } from '@swimlane/ngx-datatable';
// temporary 2000 name hack because of conflict
import { CloudAccount2000Component } from './system/cloud-accounts/cloud-account/cloud-account.component';
import { CloudAccounts2000Component } from './system/cloud-accounts/cloud-accounts.component';
import { Users2000Component } from './system/users/users.component';
import { EditCloudAccountComponent } from './system/cloud-accounts/edit-cloud-account/edit-cloud-account.component';
import { MainComponent } from './system/main/main.component';
import { HelmReposComponent } from './system/main/helm-repos/helm-repos.component';
import { NewClusterComponent } from './clusters/new-cluster/new-cluster.component';
import { ClusterComponent } from './clusters/cluster/cluster.component';
import { ClustersListComponent } from './clusters/clusters-list/clusters-list.component';
import { ListCloudAccountsComponent } from './system/cloud-accounts/list-cloud-accounts/list-cloud-accounts.component';
import { DashboardTutorialComponent } from './tutorials/dashboard-tutorial/dashboard-tutorial.component';
import { ClustersTutorialComponent } from './tutorials/clusters-tutorial/clusters-tutorial.component';
import { SystemTutorialComponent } from './tutorials/system-tutorial/system-tutorial.component';
import { AppsTutorialComponent } from './tutorials/apps-tutorial/apps-tutorial.component';
import { NewAppListComponent } from './apps/new-app-list/new-app-list.component';
import { NewAppComponent } from './apps/new-app/new-app.component';
import { LogsComponent } from './system/logs/logs.component';
import { NewKubeResourceComponent } from './kube-resources/new-kube-resource/new-kube-resource.component';
import { ToolbarComponent } from './navigation/toolbar/toolbar.component';
import { UserMenuComponent } from './navigation/user-menu/user-menu.component';
import { FooterComponent } from './shared/footer/footer.component';







@NgModule({
  declarations: [
    AppComponent,
    VolumesComponent,
    VolumeComponent,
    VolumesHeaderComponent,
    NavigationComponent,
    UsersComponent,
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
    HelmAppComponent,
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
    AppsListComponent,
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
    ClustersListComponent,
    ListCloudAccountsComponent,
    DashboardTutorialComponent,
    ClustersTutorialComponent,
    SystemTutorialComponent,
    AppsTutorialComponent,
    NewAppListComponent,
    NewAppComponent,
    LogsComponent,
    NewKubeResourceComponent,
    ToolbarComponent,
    UserMenuComponent,
    FooterComponent,
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    // Material:
    MatAutocompleteModule,
    MatButtonModule,
    MatButtonToggleModule,
    MatCardModule,
    MatCheckboxModule,
    MatChipsModule,
    MatDatepickerModule,
    MatDialogModule,
    MatDividerModule,
    MatExpansionModule,
    MatGridListModule,
    MatIconModule,
    MatInputModule,
    MatListModule,
    MatMenuModule,
    MatNativeDateModule,
    MatPaginatorModule,
    MatProgressBarModule,
    MatProgressSpinnerModule,
    MatRadioModule,
    MatRippleModule,
    MatSelectModule,
    MatSidenavModule,
    MatSliderModule,
    MatSlideToggleModule,
    MatSnackBarModule,
    MatSortModule,
    MatStepperModule,
    MatTableModule,
    MatTabsModule,
    MatToolbarModule,
    MatTooltipModule,
    CommonModule,
    ContextMenuModule.forRoot(),
    NgbModule.forRoot(),
    AppRoutingModule,
    FormsModule,
    BrowserAnimationsModule,
    SimpleNotificationsModule.forRoot(),
    ReactiveFormsModule,
    SchemaFormModule,
    ChartsModule,
    NgxPaginationModule,
    AceEditorModule,
    BrowserModule, MaterialDesignFrameworkModule,
    NgxDatatableModule,
  ],
  providers: [
    TitleCasePipe,
    // Component Services
    CloudAccountsService,
    SessionsService,
    UsersService,
    NodesService,
    LoadBalancersService,
    PodsService,
    ServicesService,
    VolumesService,
    AppsService,
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
    AuthService,
    AuthGuardService,
    // { provide: Http, useClass: AuthenticatedHttpService },
    {
      provide: HTTP_INTERCEPTORS,
      useClass: TokenInterceptor,
      multi: true
    }
  ],
  bootstrap: [AppComponent]
})

export class AppModule { }

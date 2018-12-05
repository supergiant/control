// Modules
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgModule } from '@angular/core';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { HttpClientModule } from '@angular/common/http';
import { SimpleNotificationsModule } from 'angular2-notifications';
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
import { SessionsHeaderComponent } from './sessions/sessions-header/sessions-header.component';
import { LoadBalancersHeaderComponent } from './load-balancers/load-balancers-header/load-balancers-header.component';
import { CloudAccountsHeaderComponent } from './cloud-accounts/cloud-accounts-header/cloud-accounts-header.component';
import { CloudAccountComponent } from './cloud-accounts/cloud-account/cloud-account.component';
import { LoadBalancerComponent } from './load-balancers/load-balancer/load-balancer.component';
import { SessionComponent } from './sessions/session/session.component';
import { UserComponent } from './users/user/user.component';
import { UsersHeaderComponent } from './users/users-header/users-header.component';
import { SystemModalComponent } from './shared/system-modal/system-modal.component';
import { DropdownModalComponent } from './shared/dropdown-modal/dropdown-modal.component';
import { EditModalComponent } from './shared/edit-modal/edit-modal.component';
import { LoginComponent } from './login/login.component';
import { CookiesComponent } from './shared/cookies/cookies.component';
import { Search } from './shared/search-pipe/search-pipe';
import { SessionDetailsComponent } from './sessions/session-details/session-details.component';
import { SessionsListComponent } from './sessions/sessions-list/sessions-list.component';
import { UserDetailsComponent } from './users/user-details/user-details.component';
import { UsersListComponent } from './users/users-list/users-list.component';
import { CloudAccountDetailsComponent } from './cloud-accounts/cloud-account-details/cloud-account-details.component';
import { CloudAccountsListComponent } from './cloud-accounts/cloud-accounts-list/cloud-accounts-list.component';
import { LoadBalancerDetailsComponent } from './load-balancers/load-balancer-details/load-balancer-details.component';
import { LoadBalancersListComponent } from './load-balancers/load-balancers-list/load-balancers-list.component';
import { SupergiantComponent } from './shared/supergiant/supergiant.component';
// Component Services
import { SessionsService } from './sessions/sessions.service';
import { CloudAccountsService } from './cloud-accounts/cloud-accounts.service';
import { UsersService } from './users/users.service';
import { NodesService } from './nodes/nodes.service';
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
import { AuthService } from './shared/supergiant/auth/auth.service';
import { AuthGuardService } from './shared/supergiant/auth/auth-guard.service';
import { TokenInterceptor } from './shared/supergiant/auth/token.interceptor';
import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { WINDOW_PROVIDERS } from './shared/helpers/window-providers';








// ui2000
import { SystemComponent } from './system/system.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { ClustersComponent } from './clusters/clusters.component';
import { NewCloudAccountComponent } from './system/cloud-accounts/new-cloud-account/new-cloud-account.component';
import { NgxDatatableModule } from '@swimlane/ngx-datatable';
// temporary 2000 name hack because of conflict
import { CloudAccount2000Component } from './system/cloud-accounts/cloud-account/cloud-account.component';
import { CloudAccounts2000Component } from './system/cloud-accounts/cloud-accounts.component';
import { EditCloudAccountComponent } from './system/cloud-accounts/edit-cloud-account/edit-cloud-account.component';
import { NewClusterComponent } from './clusters/new-cluster/new-cluster.component';
import { ClusterComponent } from './clusters/cluster/cluster.component';
import { ListCloudAccountsComponent } from './system/cloud-accounts/list-cloud-accounts/list-cloud-accounts.component';
import { ToolbarComponent } from './navigation/toolbar/toolbar.component';
import { UserMenuComponent } from './navigation/user-menu/user-menu.component';
import { FooterComponent } from './shared/footer/footer.component';
import { ConfirmModalComponent } from './shared/modals/confirm-modal/confirm-modal.component';
import { UsageChartComponent } from './clusters/cluster/usage-chart/usage-chart.component';
import { StoreModule } from '@ngrx/store';
import { reducers, metaReducers } from './reducers';
import { StoreDevtoolsModule } from '@ngrx/store-devtools';
import { environment } from '../environments/environment';
import { EffectsModule } from '@ngrx/effects';
import { AppEffects } from './app.effects';
import { TaskLogsComponent } from './clusters/cluster/task-logs/task-logs.component';
import { LoginGuardService } from "./shared/supergiant/auth/login-guard.service";
import { MenuModalComponent } from './navigation/user-menu/menu-modal/menu-modal.component';
import { ClusterListModalComponent } from './navigation/toolbar/cluster-list-modal/cluster-list-modal.component';
import { AddNodeComponent } from './clusters/cluster/add-node/add-node.component';
import { UsageOrbComponent } from './dashboard/usage-orb/usage-orb.component';
import { ClusterTableComponent } from './dashboard/cluster-table/cluster-table.component';
import { DeleteClusterModalComponent } from './clusters/cluster/delete-cluster-modal/delete-cluster-modal.component';
import { DeleteReleaseModalComponent } from './clusters/cluster/delete-release-modal/delete-release-modal.component';
import { SshCommandsModalComponent }   from './clusters/cluster/ssh-commands-modal/ssh-commands-modal.component';
import { KubectlConfigModalComponent } from './clusters/cluster/kubectl-config-modal/kubectl-config-modal.component';
import { ReleaseInfoModalComponent }   from './clusters/cluster/release-info-modal/release-info-modal.component';
import { CallbackPipe }                from './callback.pipe';

@NgModule({
  declarations: [
    AppComponent,
    NavigationComponent,
    UsersComponent,
    SessionsComponent,
    CloudAccountsComponent,
    LoadBalancersComponent,
    SessionsHeaderComponent,
    LoadBalancersHeaderComponent,
    CloudAccountsHeaderComponent,
    CloudAccountComponent,
    LoadBalancerComponent,
    SessionComponent,
    UserComponent,
    UsersHeaderComponent,
    NotificationsComponent,
    SystemModalComponent,
    DropdownModalComponent,
    EditModalComponent,
    LoginComponent,
    CookiesComponent,
    SessionDetailsComponent,
    SessionsListComponent,
    UserDetailsComponent,
    UsersListComponent,
    CloudAccountDetailsComponent,
    CloudAccountsListComponent,
    LoadBalancerDetailsComponent,
    LoadBalancersListComponent,
    Search,
    SupergiantComponent,
    DashboardComponent,
    ClustersComponent,
    NewCloudAccountComponent,
    EditCloudAccountComponent,
    NewClusterComponent,
    ClusterComponent,
    CloudAccount2000Component,
    CloudAccounts2000Component,
    SystemComponent,
    ListCloudAccountsComponent,
    ToolbarComponent,
    UserMenuComponent,
    FooterComponent,
    ConfirmModalComponent,
    UsageChartComponent,
    TaskLogsComponent,
    MenuModalComponent,
    ClusterListModalComponent,
    AddNodeComponent,
    UsageOrbComponent,
    ClusterTableComponent,
    DeleteClusterModalComponent,
    DeleteReleaseModalComponent,
    SshCommandsModalComponent,
    KubectlConfigModalComponent,
    ReleaseInfoModalComponent,
    CallbackPipe,
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
    StoreModule.forRoot(reducers, { metaReducers }),
    !environment.production ? StoreDevtoolsModule.instrument() : [],
    EffectsModule.forRoot([AppEffects]),
  ],
  providers: [
    WINDOW_PROVIDERS,
    TitleCasePipe,
    // Component Services
    CloudAccountsService,
    SessionsService,
    UsersService,
    NodesService,
    LoadBalancersService,
    ServicesService,
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
    LoginGuardService,
    // { provide: Http, useClass: AuthenticatedHttpService },
    {
      provide: HTTP_INTERCEPTORS,
      useClass: TokenInterceptor,
      multi: true
    }
  ],
  entryComponents: [
    ConfirmModalComponent,
    TaskLogsComponent,
    MenuModalComponent,
    ClusterListModalComponent,
    DeleteClusterModalComponent,
    DeleteReleaseModalComponent,
    SshCommandsModalComponent,
    KubectlConfigModalComponent,
    ReleaseInfoModalComponent
  ],
  bootstrap: [AppComponent],
})

export class AppModule { }

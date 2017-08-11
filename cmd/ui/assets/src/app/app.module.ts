// Modules
import { BrowserModule } from '@angular/platform-browser';
import {FormsModule, ReactiveFormsModule} from '@angular/forms'
import { NgModule } from '@angular/core';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { HttpModule } from '@angular/http';
import { SimpleNotificationsModule } from 'angular2-notifications';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import {FormlyModule, FormlyBootstrapModule} from 'ng-formly';
import { SchemaFormModule, WidgetRegistry, DefaultWidgetRegistry } from "angular2-schema-form";
import { AppRoutingModule } from './app-routing.module';

// Components
import { AppComponent } from './app.component';
import { NavigationComponent } from './navigation/navigation.component';
import { UsersComponent } from './users/users.component';
import { KubesComponent } from './kubes/kubes.component';
import { NotificationsComponent } from './shared/notifications/notifications.component';
import { KubeComponent } from './kubes/kube/kube.component';
import { KubesHeaderComponent } from './kubes/kubes-header/kubes-header.component';
import { HeaderComponent } from './shared/header/header.component';
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
import { PodComponent } from './pods/pod/pod.component';
import { ServiceComponent } from './services/service/service.component';
import { SessionComponent } from './sessions/session/session.component';
import { UserComponent } from './users/user/user.component';
import { UsersHeaderComponent } from './users/users-header/users-header.component';
import { SystemModalComponent } from './shared/system-modal/system-modal.component';
import { DropdownModalComponent } from './shared/dropdown-modal/dropdown-modal.component';
import { EditModalComponent } from './shared/edit-modal/edit-modal.component';
// Component Services
import { SessionsService } from './sessions/sessions.service';
import { CloudAccountsService } from './cloud-accounts/cloud-accounts.service';
import { KubesService } from './kubes/kubes.service';
import { Notifications } from './shared/notifications/notifications.service';
import { SystemModalService } from './shared/system-modal/system-modal.service';
import { DropdownModalService } from './shared/dropdown-modal/dropdown-modal.service';
import { EditModalService } from './shared/edit-modal/edit-modal.service';

// Supergiant API Services
import { Supergiant } from './shared/supergiant/supergiant.service';
import { UtilService} from './shared/supergiant/util/util.service';
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




@NgModule({
  declarations: [
    AppComponent,
    NavigationComponent,
    UsersComponent,
    KubesComponent,
    KubeComponent,
    KubesHeaderComponent,
    HeaderComponent,
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
    EditModalComponent
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
    FormlyModule,
    FormlyBootstrapModule,
    ReactiveFormsModule,
    SchemaFormModule
  ],
  providers: [
    // Component Services
    KubesService,
    CloudAccountsService,
    SessionsService,
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
    {provide: WidgetRegistry, useClass: DefaultWidgetRegistry},
    Notifications,
    SystemModalService,
    DropdownModalService,
    EditModalService,
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }

import { BrowserModule } from '@angular/platform-browser';
import {FormsModule, ReactiveFormsModule} from '@angular/forms'
import { NgModule } from '@angular/core';
import { AppComponent } from './app.component';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { NavigationComponent } from './navigation/navigation.component';
import { UsersComponent } from './users/users.component';
import { KubesComponent } from './kubes/kubes.component';
import { AppRoutingModule } from './app-routing.module';
import { KubeComponent } from './kubes/kube/kube.component';
import { KubesService } from './kubes/kubes.service';
import { KubeHeaderComponent } from './kubes/kube-header/kube-header.component';
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
import { CloudAccountsNewModalComponent } from './cloud-accounts/cloud-accounts-new-modal/cloud-accounts-new-modal.component';
import { CloudAccountsService } from './cloud-accounts/cloud-accounts.service';
import { CloudAccountsNewSubmitComponent } from './cloud-accounts/cloud-accounts-new-submit/cloud-accounts-new-submit.component';
import { HttpModule } from '@angular/http';
import { SimpleNotificationsModule } from 'angular2-notifications';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import {FormlyModule, FormlyBootstrapModule} from 'ng-formly';
import { SchemaFormModule, WidgetRegistry, DefaultWidgetRegistry } from "angular2-schema-form";
import { SupergiantComponent } from './shared/supergiant/supergiant.component';
import { Supergiant } from './shared/supergiant/supergiant.service';


@NgModule({
  declarations: [
    AppComponent,
    NavigationComponent,
    UsersComponent,
    KubesComponent,
    KubeComponent,
    KubeHeaderComponent,
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
    CloudAccountsNewModalComponent,
    CloudAccountsNewSubmitComponent,
    SupergiantComponent,
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
  providers: [KubesService, CloudAccountsService, {provide: WidgetRegistry, useClass: DefaultWidgetRegistry}],
  bootstrap: [AppComponent]
})
export class AppModule { }

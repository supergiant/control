import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { NodesComponent } from './nodes/nodes.component';
import { ServicesComponent } from './services/services.component';
import { AppsComponent } from './apps/apps.component';
import { LoginComponent } from './login/login.component';
import { NodeDetailsComponent } from './nodes/node-details/node-details.component';
import { ServiceDetailsComponent } from './services/service-details/service-details.component';
import { ServicesListComponent } from './services/services-list/services-list.component';
import { AppsListComponent } from './apps/apps-list/apps-list.component';
import { NewKubeResourceComponent } from './kube-resources/new-kube-resource/new-kube-resource.component';

// ui 2000 components
import { SystemComponent } from './system/system.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { ClustersComponent } from './clusters/clusters.component';
import { NewCloudAccountComponent } from './system/cloud-accounts/new-cloud-account/new-cloud-account.component';
// temporary 2000 name hack because of conflict
import { CloudAccount2000Component } from './system/cloud-accounts/cloud-account/cloud-account.component';
import { CloudAccounts2000Component } from './system/cloud-accounts/cloud-accounts.component';
import { ListCloudAccountsComponent } from './system/cloud-accounts/list-cloud-accounts/list-cloud-accounts.component';
import { Users2000Component } from './system/users/users.component';
import { EditCloudAccountComponent } from './system/cloud-accounts/edit-cloud-account/edit-cloud-account.component';
import { MainComponent } from './system/main/main.component';
import { NewClusterComponent } from './clusters/new-cluster/new-cluster.component';
import { ClusterComponent } from './clusters/cluster/cluster.component';
import { ClustersListComponent } from './clusters/clusters-list/clusters-list.component';
import { DashboardTutorialComponent } from './tutorials/dashboard-tutorial/dashboard-tutorial.component';
import { SystemTutorialComponent } from './tutorials/system-tutorial/system-tutorial.component';
import { AppsTutorialComponent } from './tutorials/apps-tutorial/apps-tutorial.component';
import { NewAppListComponent } from './apps/new-app-list/new-app-list.component';
import { NewAppComponent } from './apps/new-app/new-app.component';
import { LogsComponent } from './system/logs/logs.component';
// auth guard
import {
  AuthGuardService as AuthGuard
} from './shared/supergiant/auth/auth-guard.service';

const appRoutes: Routes = [
  { path: '', component: LoginComponent },
  {
    path: 'dashboard', component: DashboardComponent, canActivate: [AuthGuard], children: [
      { path: '', component: DashboardTutorialComponent },
    ]
  },
  {
    path: 'apps', component: AppsComponent, canActivate: [AuthGuard], children: [
      { path: '', component: AppsTutorialComponent },
      { path: '', component: AppsListComponent },
      { path: 'new', component: NewAppListComponent },
      { path: 'new/:id', component: NewAppComponent },
    ]
  },
  {
    path: 'clusters', component: ClustersComponent, canActivate: [AuthGuard], children: [
      { path: '', component: ClustersListComponent },
      { path: 'new', component: NewClusterComponent },
      { path: ':id', component: ClusterComponent }
    ]
  },
  {
    path: 'system', component: SystemComponent, canActivate: [AuthGuard], children: [
      { path: '', component: SystemTutorialComponent },
      {
        path: 'logs', component: LogsComponent, children: [
          { path: '', component: LogsComponent },
        ],
      },
      {
        path: 'cloud-accounts', component: CloudAccounts2000Component, children: [
          { path: '', component: ListCloudAccountsComponent },
          { path: 'new', component: NewCloudAccountComponent },
          { path: 'edit/:id', component: EditCloudAccountComponent },
          { path: ':id', component: CloudAccount2000Component },
        ],
      },
      {
        path: 'users', component: Users2000Component, children: [
          { path: 'new', component: NewCloudAccountComponent },
          { path: 'edit/:id', component: EditCloudAccountComponent },
          { path: ':id', component: CloudAccount2000Component },
        ]
      },
      { path: 'main', component: MainComponent },
      { path: '', component: MainComponent },
    ]
  },
  {
    path: 'nodes', component: NodesComponent, canActivate: [AuthGuard], children: [
      { path: ':id', component: NodeDetailsComponent }
    ]
  },
  {
    path: 'resource/new/:id', component: NewKubeResourceComponent, canActivate: [AuthGuard]
  },
  {
    path: 'services', component: ServicesComponent, canActivate: [AuthGuard], children: [
      { path: '', component: ServicesListComponent },
      { path: ':id', component: ServiceDetailsComponent }
    ]
  },
];

@NgModule({
  imports: [RouterModule.forRoot(appRoutes)],
  exports: [RouterModule],
  providers: [AuthGuard]
})
export class AppRoutingModule {

}

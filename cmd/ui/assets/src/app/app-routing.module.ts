import { NgModule, Injectable } from '@angular/core';
import { Routes, Router, RouterModule, CanActivate } from '@angular/router';
import { UsersComponent } from './users/users.component';
import { CloudAccountsComponent } from './cloud-accounts/cloud-accounts.component';
import { NodesComponent } from './nodes/nodes.component';
import { ServicesComponent } from './services/services.component';
import { SessionsComponent } from './sessions/sessions.component';
import { PodsComponent } from './pods/pods.component';
import { AppsComponent } from './apps/apps.component';
import { LoginComponent } from './login/login.component';
import { VolumesComponent } from './volumes/volumes.component';
import { LoadBalancersComponent } from './load-balancers/load-balancers.component';
import { Supergiant } from './shared/supergiant/supergiant.service';
import { Observable } from 'rxjs/Observable';
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
import { HelmReposComponent } from './system/main/helm-repos/helm-repos.component';
import { NewClusterComponent } from './clusters/new-cluster/new-cluster.component';
import { ClusterComponent } from './clusters/cluster/cluster.component';
import { ClustersListComponent } from './clusters/clusters-list/clusters-list.component';
import { DashboardTutorialComponent } from './tutorials/dashboard-tutorial/dashboard-tutorial.component';
import { ClustersTutorialComponent } from './tutorials/clusters-tutorial/clusters-tutorial.component';
import { SystemTutorialComponent } from './tutorials/system-tutorial/system-tutorial.component';
import { AppsTutorialComponent } from './tutorials/apps-tutorial/apps-tutorial.component';
import { NewAppListComponent } from './apps/new-app-list/new-app-list.component';
import { NewAppComponent } from './apps/new-app/new-app.component';
import { LogsComponent } from './system/logs/logs.component';

@Injectable()
export class AuthGuard implements CanActivate {

  constructor(
    private router: Router,
    private supergiant: Supergiant,
    private loginComponent: LoginComponent,
  ) { }

  canActivate(): Observable<boolean> | boolean {
    return this.loginComponent.validateUser().map((res) => {
      if (res) { return true; }
    }).catch(() => {
      this.router.navigate(['']);
      return Observable.of(false);
    });
  }

  handleError() {
    // this.router.navigate(['/login']);
    return Observable.of(false);
  }
}
const appRoutes: Routes = [
  { path: '', component: LoginComponent },
  {
    path: 'dashboard', component: DashboardComponent, canActivate: [AuthGuard], children: [
      { path: '', component: DashboardTutorialComponent, outlet: 'tutorial' },
    ]
  },
  {
    path: 'apps', component: AppsComponent, canActivate: [AuthGuard], children: [
      { path: '', component: AppsTutorialComponent, outlet: 'tutorial' },
      { path: '', component: AppsListComponent },
      { path: 'new', component: NewAppListComponent },
      { path: 'new/:id', component: NewAppComponent },
      // { path: ':id', component: AppDetailsComponent },
      // { path: 'deployable/:id', component: DeploymentDetailsComponent },
    ]
  },
  {
    path: 'clusters', component: ClustersComponent, canActivate: [AuthGuard], children: [
      { path: '', component: ClustersTutorialComponent, outlet: 'tutorial' },
      { path: '', component: ClustersListComponent },
      { path: 'new', component: NewClusterComponent },
      { path: ':id', component: ClusterComponent }
    ]
  },
  {
    path: 'system', component: SystemComponent, canActivate: [AuthGuard], children: [
      { path: '', component: SystemTutorialComponent, outlet: 'tutorial' },
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
  // {
  //   path: 'kubes', component: KubesComponent, canActivate: [AuthGuard], children: [
  //     { path: '', component: KubesListComponent },
  //     { path: ':id', component: KubeDetailsComponent }
  //   ]
  // },
  // {
  //   path: 'users', component: UsersComponent, canActivate: [AuthGuard], children: [
  //     { path: '', component: UsersListComponent },
  //     { path: ':id', component: UserDetailsComponent }
  //   ]
  // },
  // {
  //   path: 'cloud-accounts', component: CloudAccountsComponent, canActivate: [AuthGuard], children: [
  //     { path: '', component: CloudAccountsListComponent },
  //     { path: ':id', component: CloudAccountDetailsComponent }
  //   ]
  // },
  {
    path: 'nodes', component: NodesComponent, canActivate: [AuthGuard], children: [
      { path: '', component: NodesListComponent },
      { path: ':id', component: NodeDetailsComponent }
    ]
  },
  {
    path: 'pods', component: PodsComponent, canActivate: [AuthGuard], children: [
      { path: '', component: PodsListComponent },
      { path: ':id', component: PodDetailsComponent }
    ]
  },
  {
    path: 'resource/new/:id', component: NewKubeResourceComponent, canActivate: [AuthGuard]
  }
  // { path: 'volumes', component: VolumesComponent, canActivate: [AuthGuard] },
  // {
  //   path: 'services', component: ServicesComponent, canActivate: [AuthGuard], children: [
  //     { path: '', component: ServicesListComponent },
  //     { path: ':id', component: ServiceDetailsComponent }
  //   ]
  // },
  // {
  //   path: 'sessions', component: SessionsComponent, canActivate: [AuthGuard], children: [
  //     { path: '', component: SessionsListComponent },
  //     { path: ':id', component: SessionDetailsComponent }
  //   ]
  // },
  // {
  //   path: 'load-balancers', component: LoadBalancersComponent, canActivate: [AuthGuard], children: [
  //     { path: '', component: LoadBalancersListComponent },
  //     { path: ':id', component: LoadBalancerDetailsComponent }
  //   ]
  // },
  // { path: 'login', component: LoginComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(appRoutes)],
  exports: [RouterModule],
  providers: [AuthGuard]
})
export class AppRoutingModule {

}

import { NgModule, Injectable } from '@angular/core';
import { Routes, Router, RouterModule, CanActivate } from '@angular/router';
import { KubesComponent } from './kubes/kubes.component';
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
      this.router.navigate(['/login']);
      return Observable.of(false);
    });
  }

  handleError() {
    // this.router.navigate(['/login']);
    return Observable.of(false);
  }
}
const appRoutes: Routes = [
  { path: '', redirectTo: '/login', pathMatch: 'full' },
  {
    path: 'kubes', component: KubesComponent, canActivate: [AuthGuard], children: [
      { path: '', component: KubesListComponent },
      { path: ':id', component: KubeDetailsComponent }
    ]
  },
  {
    path: 'users', component: UsersComponent, canActivate: [AuthGuard], children: [
      { path: '', component: UsersListComponent },
      { path: ':id', component: UserDetailsComponent }
    ]
  },
  {
    path: 'cloud-accounts', component: CloudAccountsComponent, canActivate: [AuthGuard], children: [
      { path: '', component: CloudAccountsListComponent },
      { path: ':id', component: CloudAccountDetailsComponent }
    ]
  },
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
  { path: 'apps', component: AppsComponent, canActivate: [AuthGuard] },
  { path: 'volumes', component: VolumesComponent, canActivate: [AuthGuard] },
  {
    path: 'services', component: ServicesComponent, canActivate: [AuthGuard], children: [
      { path: '', component: ServicesListComponent },
      { path: ':id', component: ServiceDetailsComponent }
    ]
  },
  {
    path: 'sessions', component: SessionsComponent, canActivate: [AuthGuard], children: [
      { path: '', component: SessionsListComponent },
      { path: ':id', component: SessionDetailsComponent }
    ]
  },
  {
    path: 'load-balancers', component: LoadBalancersComponent, canActivate: [AuthGuard], children: [
      { path: '', component: LoadBalancersListComponent },
      { path: ':id', component: LoadBalancerDetailsComponent }
    ]
  },
  { path: 'login', component: LoginComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(appRoutes)],
  exports: [RouterModule],
  providers: [AuthGuard]
})
export class AppRoutingModule {

}

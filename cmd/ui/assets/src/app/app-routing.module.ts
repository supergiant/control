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
      { path: '', component: KubesListComponent }
    ]
  },
  {
    path: 'kubes/:id', component: KubeDetailsComponent, canActivate: [AuthGuard], children: [
      { path: 'details', component: KubeDetailsComponent }
    ]
  },
  { path: 'users', component: UsersComponent, canActivate: [AuthGuard] },
  { path: 'cloud-accounts', component: CloudAccountsComponent, canActivate: [AuthGuard] },
  { path: 'nodes', component: NodesComponent, canActivate: [AuthGuard] },
  { path: 'pods', component: PodsComponent, canActivate: [AuthGuard] },
  { path: 'apps', component: AppsComponent, canActivate: [AuthGuard] },
  { path: 'volumes', component: VolumesComponent, canActivate: [AuthGuard] },
  { path: 'services', component: ServicesComponent, canActivate: [AuthGuard] },
  { path: 'sessions', component: SessionsComponent, canActivate: [AuthGuard] },
  { path: 'load-balancers', component: LoadBalancersComponent, canActivate: [AuthGuard] },
  { path: 'login', component: LoginComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(appRoutes)],
  exports: [RouterModule],
  providers: [AuthGuard]
})
export class AppRoutingModule {

}

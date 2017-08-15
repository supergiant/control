import { NgModule, Injectable } from '@angular/core';
import { Routes, Router, RouterModule, CanActivate } from '@angular/router';

import { KubesComponent } from './kubes/kubes.component';
import { UsersComponent } from './users/users.component';
import { CloudAccountsComponent } from './cloud-accounts/cloud-accounts.component';
import { NodesComponent } from './nodes/nodes.component';
import { ServicesComponent } from './services/services.component';
import { SessionsComponent } from './sessions/sessions.component';
import { PodsComponent } from './pods/pods.component';
import { LoginComponent } from './login/login.component';
import { VolumesComponent } from './volumes/volumes.component';
import { LoadBalancersComponent } from './load-balancers/load-balancers.component';
import { Supergiant } from './shared/supergiant/supergiant.service'

@Injectable()
export class AuthGuard implements CanActivate {

    constructor(
      private router: Router,
      private supergiant: Supergiant,
    ) { }

    canActivate() {
      if (this.supergiant.loginSuccess) {return true}
        this.router.navigate(['/login']);
        return false;
    }
}

const appRoutes: Routes = [
  { path: '', redirectTo: '/kubes', pathMatch: 'full' },
  { path: 'kubes', component: KubesComponent, canActivate: [AuthGuard], children: [
    { path: '', component: KubesComponent },
    { path: 'new', component: KubesComponent },
    { path: ':id', component: KubesComponent}
  ] },
  { path: 'users', component: UsersComponent, canActivate: [AuthGuard] },
  { path: 'cloud-accounts', component: CloudAccountsComponent, canActivate: [AuthGuard] },
  { path: 'nodes', component: NodesComponent, canActivate: [AuthGuard] },
  { path: 'pods', component: PodsComponent, canActivate: [AuthGuard] },
  { path: 'volumes', component: VolumesComponent, canActivate: [AuthGuard] },
  { path: 'services', component: ServicesComponent, canActivate: [AuthGuard] },
  { path: 'sessions', component: SessionsComponent, canActivate: [AuthGuard] },
  { path: 'load-balancers', component: LoadBalancersComponent, canActivate: [AuthGuard] },
  { path: 'login', component: LoginComponent},
];

@NgModule({
  imports: [RouterModule.forRoot(appRoutes)],
  exports: [RouterModule],
  providers: [AuthGuard]
})
export class AppRoutingModule {

}

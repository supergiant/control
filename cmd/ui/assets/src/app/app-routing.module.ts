import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { KubesComponent } from './kubes/kubes.component';
import { UsersComponent } from './users/users.component';
import { CloudAccountsComponent } from './cloud-accounts/cloud-accounts.component';
import { NodesComponent } from './nodes/nodes.component';
import { ServicesComponent } from './services/services.component';
import { SessionsComponent } from './sessions/sessions.component';
import { PodsComponent } from './pods/pods.component';
import { LoadBalancersComponent } from './load-balancers/load-balancers.component';

const appRoutes: Routes = [
  { path: '', redirectTo: '/kubes', pathMatch: 'full' },
  { path: 'kubes', component: KubesComponent, children: [
    { path: '', component: KubesComponent },
    { path: 'new', component: KubesComponent },
    { path: ':id', component: KubesComponent}
  ] },
  { path: 'users', component: UsersComponent },
  { path: 'cloud-accounts', component: CloudAccountsComponent },
  { path: 'nodes', component: NodesComponent },
  { path: 'pods', component: PodsComponent },
  { path: 'services', component: ServicesComponent },
  { path: 'sessions', component: SessionsComponent },
  { path: 'load-balancers', component: LoadBalancersComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(appRoutes)],
  exports: [RouterModule]
})
export class AppRoutingModule {

}

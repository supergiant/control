import { NgModule }                 from '@angular/core';
import { Routes, RouterModule }     from '@angular/router';
import { AddNodeComponent }         from './clusters/cluster/add-node/add-node.component';
import { LoginComponent }           from './login/login.component';

// ui 2000 components
import { SystemComponent } from './system/system.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { ClustersComponent } from './clusters/clusters.component';
import { NewCloudAccountComponent } from './system/cloud-accounts/new-cloud-account/new-cloud-account.component';
// temporary 2000 name hack because of conflict
import { CloudAccount2000Component } from './system/cloud-accounts/cloud-account/cloud-account.component';
import { CloudAccounts2000Component } from './system/cloud-accounts/cloud-accounts.component';
import { ListCloudAccountsComponent } from './system/cloud-accounts/list-cloud-accounts/list-cloud-accounts.component';
import { EditCloudAccountComponent } from './system/cloud-accounts/edit-cloud-account/edit-cloud-account.component';
import { NewClusterComponent } from './clusters/new-cluster/new-cluster.component';
import { ClusterComponent } from './clusters/cluster/cluster.component';

// auth guard
import {
  AuthGuardService as AuthGuard
} from './shared/supergiant/auth/auth-guard.service';
import {LoginGuardService} from './shared/supergiant/auth/login-guard.service';

const appRoutes: Routes = [
  {
    path: '', component: LoginComponent,
    canActivate: [LoginGuardService]
  },
  {
    path: 'dashboard', component: DashboardComponent, canActivate: [AuthGuard]  },
  {
    path: 'catalog',
    loadChildren: 'app/features/app-store/app-store.module#AppStoreModule'
  },
  {
    path: 'clusters', component: ClustersComponent, canActivate: [AuthGuard], children: [
      { path: 'new', component: NewClusterComponent },
      { path: ':id', component: ClusterComponent },
      { path: ':id/add-node', component: AddNodeComponent },
    ]
  },
  {
    path: 'system', component: SystemComponent, canActivate: [AuthGuard], children: [
      {
        path: 'cloud-accounts', component: CloudAccounts2000Component, children: [
          { path: '', component: ListCloudAccountsComponent },
          { path: 'new', component: NewCloudAccountComponent },
          { path: 'edit/:id', component: EditCloudAccountComponent },
          { path: ':id', component: CloudAccount2000Component },
        ]
      }
    ]
  },
];

@NgModule({
  imports: [RouterModule.forRoot(appRoutes, { scrollPositionRestoration: 'enabled' })],
  exports: [RouterModule],
  providers: [AuthGuard]
})
export class AppRoutingModule {

}

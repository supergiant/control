import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { KubesComponent } from './kubes/kubes.component';
import { UsersComponent } from './users/users.component';

const appRoutes: Routes = [
  { path: '', redirectTo: '/kubes', pathMatch: 'full' },
  { path: 'kubes', component: KubesComponent, children: [
    { path: '', component: KubesComponent },
    { path: 'new', component: KubesComponent },
    { path: ':id', component: KubesComponent}
  ] },
  { path: 'users', component: UsersComponent }
];

@NgModule({
  imports: [RouterModule.forRoot(appRoutes)],
  exports: [RouterModule]
})
export class AppRoutingModule {

}

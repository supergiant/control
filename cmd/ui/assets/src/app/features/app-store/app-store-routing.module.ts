import { NgModule }                from '@angular/core';
import { Routes, RouterModule }    from '@angular/router';
import { AppStoreComponent }       from './app-store.component';
import { AppDetailsComponent }     from "./app-details/app-details.component";
import { AppsListComponent }       from "./apps-list/apps-list.component";

const routes: Routes = [
  {
    path: '',
    component: AppStoreComponent,
    children: [
      {
        path: '',
        pathMatch: 'full',
        redirectTo: 'supergiant'
      },
      {
        path: ':repo',
        component: AppsListComponent,
      },
      {
        path: ':repo/details/:chart',
        component: AppDetailsComponent
      }
    ]
  },

];

@NgModule({
  imports: [ RouterModule.forChild(routes) ],
  exports: [ RouterModule ],
})
export class AppStoreRoutingModule {
}

import { NgModule }                from '@angular/core';
import { Routes, RouterModule }    from '@angular/router';
import { AppStoreComponent }       from './app-store.component';
import { AppsSupergiantComponent } from './apps-supergiant/apps-supergiant.component';
import { AppsVerifiedComponent }   from './apps-verified/apps-verified.component';
import { AppsOtherComponent }      from './apps-other/apps-other.component';
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
        path: 'verified',
        component: AppsVerifiedComponent,
      },
      {
        path: 'others',
        component: AppsOtherComponent,
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

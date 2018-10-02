import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { AppStoreComponent } from './app-store.component';
import { AppListComponent } from './app-list/app-list.component';

const routes: Routes = [
  {
    path: '',
    redirectTo: 'list',
    pathMatch: 'full'
  },
  {
    path: 'list',
    component: AppStoreComponent,
    children: [
      {
        path: '',
        component: AppListComponent,
      },
    ]
  }
];

@NgModule({
  imports: [ RouterModule.forChild(routes) ],
  exports: [ RouterModule ],
})
export class AppStoreRoutingModule {
}

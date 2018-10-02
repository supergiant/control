import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { AppStoreRoutingModule } from './app-store-routing.module';
import { AppStoreComponent } from './app-store.component';
import { AppListComponent } from './app-list/app-list.component';

@NgModule({
  imports: [
    CommonModule,
    AppStoreRoutingModule
  ],
  declarations: [AppStoreComponent, AppListComponent]
})
export class AppStoreModule { }

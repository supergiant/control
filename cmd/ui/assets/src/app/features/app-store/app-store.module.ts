import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { AppStoreRoutingModule } from './app-store-routing.module';
import { AppStoreComponent } from './app-store.component';

@NgModule({
  imports: [
    CommonModule,
    AppStoreRoutingModule
  ],
  declarations: [AppStoreComponent]
})
export class AppStoreModule { }

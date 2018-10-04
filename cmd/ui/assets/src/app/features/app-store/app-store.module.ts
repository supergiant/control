import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { AppStoreRoutingModule } from './app-store-routing.module';
import { AppStoreComponent } from './app-store.component';
import { AppsSupergiantComponent } from './apps-supergiant/apps-supergiant.component';
import { AppsVerifiedComponent } from './apps-verified/apps-verified.component';
import { AppsOtherComponent } from './apps-other/apps-other.component';
import { AppsAddComponent } from './apps-add/apps-add.component';
import { MatCardModule } from '@angular/material';

@NgModule({
  imports: [
    CommonModule,
    AppStoreRoutingModule,
    MatCardModule
  ],
  declarations: [AppStoreComponent, AppsSupergiantComponent, AppsVerifiedComponent, AppsOtherComponent, AppsAddComponent]
})
export class AppStoreModule { }

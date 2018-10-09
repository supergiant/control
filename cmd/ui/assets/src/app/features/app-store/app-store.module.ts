import { NgModule }     from '@angular/core';
import { CommonModule } from '@angular/common';

import { AppStoreRoutingModule }   from './app-store-routing.module';
import { AppStoreComponent }       from './app-store.component';
import { AppsSupergiantComponent } from './apps-supergiant/apps-supergiant.component';
import { AppsVerifiedComponent }   from './apps-verified/apps-verified.component';
import { AppsOtherComponent }      from './apps-other/apps-other.component';
import { AppsAddComponent }        from './apps-add/apps-add.component';
import { MatCardModule }           from '@angular/material';
import { StoreModule }             from '@ngrx/store';
import * as fromApps               from '../apps/apps.reducer';
import { EffectsModule }           from '@ngrx/effects';
import { AppsEffects }             from '../apps/apps.effects';
import { AppsListComponent }       from './apps-list/apps-list.component';

@NgModule({
  imports: [
    CommonModule,
    AppStoreRoutingModule,
    MatCardModule,
    StoreModule.forFeature('apps', fromApps.reducer),
    EffectsModule.forFeature([ AppsEffects ])
  ],
  declarations: [
    AppStoreComponent,
    AppsSupergiantComponent,
    AppsVerifiedComponent,
    AppsOtherComponent,
    AppsAddComponent,
    AppsListComponent
  ]
})
export class AppStoreModule {
}

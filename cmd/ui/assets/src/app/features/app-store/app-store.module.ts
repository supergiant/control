import { NgModule }                         from '@angular/core';
import { CommonModule }                     from '@angular/common';
import { AppStoreRoutingModule }            from './app-store-routing.module';
import { AppStoreComponent }                from './app-store.component';
import { AppsAddComponent }                 from './apps-add/apps-add.component';
import {
  MatCardModule,
  MatDialogModule,
  MatFormFieldModule,
  MatInputModule,
  MatOptionModule,
  MatPaginatorModule,
  MatProgressSpinnerModule,
  MatSelectModule,
  MatTabsModule,
  MatTooltipModule,
} from '@angular/material';
import { StoreModule }                      from '@ngrx/store';
import * as fromApps                        from '../apps/apps.reducer';
import { EffectsModule }                    from '@ngrx/effects';
import { AppsEffects }                      from '../apps/apps.effects';
import { AppsListComponent }                from './apps-list/apps-list.component';
import { AppDetailsComponent }              from './app-details/app-details.component';
import { BreadcrumbsComponent }             from './breadcrumbs/breadcrumbs.component';
import { DeployComponent }                  from './app-details/deploy/deploy.component';
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { ConfigureComponent }               from './app-details/confure/configure.component';
import { MarkdownModule }                   from "ngx-markdown";
import { AceEditorModule }                  from "ng2-ace-editor";
import { RemoveRepoDialogComponent }        from './apps-list/remove-repo-dialog/remove-repo-dialog.component';
import { SortReposPipe } from './sort-repos.pipe';

@NgModule({
  imports: [
    CommonModule,
    AppStoreRoutingModule,
    MatCardModule,
    StoreModule.forFeature('apps', fromApps.reducer),
    EffectsModule.forFeature([AppsEffects]),
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatOptionModule,
    MatDialogModule,
    FormsModule,
    MatTabsModule,
    MatProgressSpinnerModule,
    MarkdownModule.forRoot(),
    AceEditorModule,
    MatPaginatorModule,
    MatTooltipModule,
  ],
  declarations: [
    AppStoreComponent,
    AppsAddComponent,
    AppsListComponent,
    AppDetailsComponent,
    BreadcrumbsComponent,
    DeployComponent,
    ConfigureComponent,
    RemoveRepoDialogComponent,
    SortReposPipe,
  ],
  entryComponents: [
    DeployComponent,
    AppsAddComponent,
    ConfigureComponent,
    RemoveRepoDialogComponent,
  ],
})
export class AppStoreModule {
}

import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { AppComponent } from './app.component';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { NavigationComponent } from './navigation/navigation.component';
import { UsersComponent } from './users/users.component';
import { KubesComponent } from './kubes/kubes.component';
import { AppRoutingModule } from './app-routing.module';
import { KubeComponent } from './kubes/kube/kube.component';
import { KubesService } from './kubes/kubes.service';
import { KubeHeaderComponent } from './kubes/kube-header/kube-header.component';
import { HeaderComponent } from './shared/header/header.component';


@NgModule({
  declarations: [
    AppComponent,
    NavigationComponent,
    UsersComponent,
    KubesComponent,
    KubeComponent,
    KubeHeaderComponent,
    HeaderComponent,
  ],
  imports: [
    BrowserModule,
    NgbModule.forRoot(),
    AppRoutingModule
  ],
  providers: [KubesService],
  bootstrap: [AppComponent]
})
export class AppModule { }

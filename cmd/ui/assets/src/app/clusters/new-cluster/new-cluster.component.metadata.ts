import { NO_ERRORS_SCHEMA } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatFormFieldModule, MatInputModule, MatSelectModule, MatStepperModule } from '@angular/material';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { RouterTestingModule } from '@angular/router/testing';
import { of } from 'rxjs';
import { CallbackPipe } from '../../callback.pipe';
import { Notifications } from '../../shared/notifications/notifications.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { NewClusterComponent } from './new-cluster.component';
import { CLOUD_ACCOUNTS_MOCK, CLUSTERS_LIST_MOCK } from './new-cluster.mocks';


class SupergiantStub {
  get Kubes() {
    return new KubesStub();
  }

  get CloudAccounts() {
    return new CloudAccountsStub();
  }
}

class KubesStub {

  get() {
    return of(CLUSTERS_LIST_MOCK);
  }

  create() {
  }
}

class CloudAccountsStub {
  get() {
    return of(CLOUD_ACCOUNTS_MOCK);
  }
}

export const NEW_CLUSTER_MODULE_METADATA = {
  schemas: [
    NO_ERRORS_SCHEMA,
  ],
  declarations: [
    NewClusterComponent,
    CallbackPipe,
  ],
  imports: [
    RouterTestingModule.withRoutes([]),
    ReactiveFormsModule,
    FormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    NoopAnimationsModule,
    MatStepperModule,
  ],
  providers: [
    {
      provide: Supergiant, useClass: SupergiantStub,
    },
    {
      provide: Notifications, useValue: { display: _ => _ },
    },
  ],
};

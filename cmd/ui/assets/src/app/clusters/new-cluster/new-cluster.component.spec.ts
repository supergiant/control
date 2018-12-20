import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA }                 from '@angular/core';

import { NewClusterComponent }                                 from './new-cluster.component';
import { CallbackPipe }                                        from 'app/callback.pipe';
import { Supergiant }                                          from 'app/shared/supergiant/supergiant.service';
import { Notifications }                                       from 'app/shared/notifications/notifications.service';
import { RouterTestingModule }                                 from '@angular/router/testing';
import { ReactiveFormsModule, FormsModule }                    from '@angular/forms';
import { of }                                                  from 'rxjs';
import { CLOUD_ACCOUNTS_MOCK }                                 from 'app/clusters/new-cluster/new-cluster.mocks';
import { MatFormFieldModule, MatInputModule, MatSelectModule } from '@angular/material';
import { NoopAnimationsModule }                                from '@angular/platform-browser/animations';

// TODO: UNIT TESTING IS REQUIRED
describe('NewClusterComponent', () => {
  let component: NewClusterComponent;
  let fixture: ComponentFixture<NewClusterComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
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
      ],
      providers: [
        {
          provide: Supergiant, useClass: SupergiantStub,
        },
        {
          provide: Notifications, useValue: { display: _ => _ },
        },
      ],
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NewClusterComponent);
    component = fixture.componentInstance;

  });

  it('should create', () => {
    const spy = spyOn(component, 'getClusters').and.callFake(_ => _);
    fixture.detectChanges();

    expect(spy).toHaveBeenCalled();
    expect(component).toBeTruthy();
  });
});


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
  }

  create() {
  }
}

class CloudAccountsStub {
  get() {
    return of(CLOUD_ACCOUNTS_MOCK);
  }
}

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA, DebugElement }   from '@angular/core';

import { NewClusterComponent }                                                   from './new-cluster.component';
import { CallbackPipe }                                                          from 'app/callback.pipe';
import { Supergiant }                                                            from 'app/shared/supergiant/supergiant.service';
import { Notifications }                                                         from 'app/shared/notifications/notifications.service';
import { RouterTestingModule }                                                   from '@angular/router/testing';
import { ReactiveFormsModule, FormsModule }                                      from '@angular/forms';
import { of }                                                                    from 'rxjs';
import {
  CLOUD_ACCOUNTS_MOCK,
  CLUSTERS_LIST_MOCK,
}                                                                                from 'app/clusters/new-cluster/new-cluster.mocks';
import { MatFormFieldModule, MatInputModule, MatSelectModule, MatStepperModule } from '@angular/material';
import { NoopAnimationsModule }                                                  from '@angular/platform-browser/animations';
import { DEFAULT_MACHINE_SET }                                                   from 'app/clusters/new-cluster/new-cluster.component.config';
import { CLUSTER_OPTIONS }                                                       from 'app/clusters/new-cluster/cluster-options.config';

// TODO: UNIT TESTING IS REQUIRED
describe('NewClusterComponent', () => {
  let component: NewClusterComponent;
  let fixture: ComponentFixture<NewClusterComponent>;

  let getClustersSpy;

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
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NewClusterComponent);
    component = fixture.componentInstance;
    getClustersSpy = spyOn(component, 'getClusters').and.callThrough();
    fixture.detectChanges();

  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should get clusters list on init', () => {
    expect(getClustersSpy).toHaveBeenCalled();
  });

  it('should set default machines list', () => {
    expect(component.machines).toEqual(DEFAULT_MACHINE_SET);
  });

  it('should set default cluster options', () => {
    expect(component.clusterOptions).toEqual(CLUSTER_OPTIONS);
  });

  it('should NOT be in provisioning mode by default', () => {
    expect(component.provisioning).toEqual(false);
  });

  describe('STEPS', () => {
    let stepper;

    describe('STEP 1: Name and Cloud account', () => {
      beforeEach(() => {
        stepper = fixture.componentInstance.stepper;
      });

      it('should NOT proceed to the next step until the form is filled', () => {
        const stepperSpy = spyOn(stepper, 'next').and.callThrough();
        component.nextStep();
        fixture.detectChanges();
        expect(stepperSpy).not.toHaveBeenCalled();
      });

      it('should PROCEED to the next step until the form is filled', () => {
        const debugEl: DebugElement = fixture.debugElement;
        const nextStepSpy = spyOn(stepper, 'next').and.callThrough();
        const controls = component.nameAndCloudAccountForm.controls;

        controls['name'].setValue('FOO');
        controls['cloudAccount'].setValue('BAR');

        debugEl
          .nativeElement
          .querySelector('button[matStepperNext]')
          .click();

        fixture.detectChanges();

        expect(nextStepSpy).toHaveBeenCalled();
      });
    });
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

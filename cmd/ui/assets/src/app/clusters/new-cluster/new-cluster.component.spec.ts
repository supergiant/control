import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { DebugElement } from '@angular/core';

import { NewClusterComponent } from './new-cluster.component';
import { DEFAULT_MACHINE_SET } from 'app/clusters/new-cluster/new-cluster.component.config';
import { CLUSTER_OPTIONS } from 'app/clusters/new-cluster/cluster-options.config';
import { NEW_CLUSTER_MODULE_METADATA } from './new-cluster.component.metadata';



// TODO: UNIT TESTING IS REQUIRED
describe('NewClusterComponent', () => {
  let component: NewClusterComponent;
  let fixture: ComponentFixture<NewClusterComponent>;

  let getClustersSpy;

  beforeEach(async(() => {
    TestBed.configureTestingModule(NEW_CLUSTER_MODULE_METADATA)
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


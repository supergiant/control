import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { DebugElement } from '@angular/core';
import { MatHorizontalStepper } from '@angular/material';

import { NewClusterComponent, StepIndexes } from './new-cluster.component';
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

  function clickBtnNext() {
    const debugEl: DebugElement = fixture.debugElement;
    debugEl
      .nativeElement
      .querySelector('button[matStepperNext]')
      .click();
  }

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

  function fillAccountDetails() {
    const controls = component.nameAndCloudAccountForm.controls;
    controls['name'].setValue('FOO');
    controls['cloudAccount'].setValue('BAR');
  }

  function goToStep2() {
    fillAccountDetails();
    clickBtnNext();
  }

  function goToStep3() {
    goToStep2();
    clickBtnNext();
  }

  describe('STEPS', () => {
    let stepper: MatHorizontalStepper;

    beforeEach(() => {
      stepper = fixture.componentInstance.stepper;
    });


    describe('STEP 1: Name and Cloud account', () => {
      it('should NOT proceed to the next step until the form is filled', () => {
        const stepperSpy = spyOn(stepper, 'next').and.callThrough();
        component.nextStep();
        fixture.detectChanges();
        expect(stepperSpy).not.toHaveBeenCalled();
      });

      it('should PROCEED to the next step until the form is filled', () => {
        const nextStepSpy = spyOn(stepper, 'next').and.callThrough();
        goToStep2();

        fixture.detectChanges();

        expect(nextStepSpy).toHaveBeenCalled();
      });
    });

    describe('STEP 2: Cluster Configuration', () => {
      beforeEach(() => {
        goToStep2();
      });

      it('should be on the Cluster Config tab', () => {
        expect(stepper.selectedIndex).toEqual(StepIndexes.ClusterConfig);
      });

      //  TODO: check all the fields set to default values
    });

    describe('STEP 3: Cluster Configuration', () => {
      beforeEach(() => {
        goToStep3();
      });

      it('should be on the Provider config tab', () => {
        expect(stepper.selectedIndex).toEqual(StepIndexes.ProvideConfig);
      });

      // TODO
      xit('should NOT PROCEED until machine config is valid', () => {
        expect(stepper.selectedIndex).toEqual(StepIndexes.ProvideConfig);
      });
    });
  });
});


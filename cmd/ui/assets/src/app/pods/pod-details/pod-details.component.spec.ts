import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';

import { PodDetailsComponent } from './pod-details.component';

describe('PodDetailsComponent', () => {
  let component: PodDetailsComponent;
  let fixture: ComponentFixture<PodDetailsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
    schemas: [NO_ERRORS_SCHEMA],
      declarations: [PodDetailsComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PodDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  xit('should be created', () => {
    expect(component).toBeTruthy();
  });
});

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';

import { PodsHeaderComponent } from './pods-header.component';

describe('PodsHeaderComponent', () => {
  let component: PodsHeaderComponent;
  let fixture: ComponentFixture<PodsHeaderComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
    schemas: [NO_ERRORS_SCHEMA],
      declarations: [ PodsHeaderComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PodsHeaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  xit('should be created', () => {
    expect(component).toBeTruthy();
  });
});

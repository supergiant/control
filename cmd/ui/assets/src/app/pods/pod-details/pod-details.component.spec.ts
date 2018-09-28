import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PodDetailsComponent } from './pod-details.component';

describe('PodDetailsComponent', () => {
  let component: PodDetailsComponent;
  let fixture: ComponentFixture<PodDetailsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
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

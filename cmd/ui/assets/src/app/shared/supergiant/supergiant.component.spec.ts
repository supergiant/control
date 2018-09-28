import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SupergiantComponent } from './supergiant.component';

describe('SupergiantComponent', () => {
  let component: SupergiantComponent;
  let fixture: ComponentFixture<SupergiantComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SupergiantComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SupergiantComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  xit('should be created', () => {
    expect(component).toBeTruthy();
  });
});

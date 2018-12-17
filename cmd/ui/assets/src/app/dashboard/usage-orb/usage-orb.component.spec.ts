import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { UsageOrbComponent } from './usage-orb.component';

describe('UsageOrbComponent', () => {
  let component: UsageOrbComponent;
  let fixture: ComponentFixture<UsageOrbComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ UsageOrbComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(UsageOrbComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  xit('should create', () => {
    expect(component).toBeTruthy();
  });
});

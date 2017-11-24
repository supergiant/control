import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DashboardTutorialComponent } from './dashboard-tutorial.component';

describe('DashboardTutorialComponent', () => {
  let component: DashboardTutorialComponent;
  let fixture: ComponentFixture<DashboardTutorialComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DashboardTutorialComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DashboardTutorialComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

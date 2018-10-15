import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AppsOtherComponent } from './apps-other.component';

describe('AppsOtherComponent', () => {
  let component: AppsOtherComponent;
  let fixture: ComponentFixture<AppsOtherComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AppsOtherComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AppsOtherComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AppsVerifiedComponent } from './apps-verified.component';

describe('AppsVerifiedComponent', () => {
  let component: AppsVerifiedComponent;
  let fixture: ComponentFixture<AppsVerifiedComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AppsVerifiedComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AppsVerifiedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

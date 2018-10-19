import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AppsSupergiantComponent } from './apps-supergiant.component';

describe('AppsSupergiantComponent', () => {
  let component: AppsSupergiantComponent;
  let fixture: ComponentFixture<AppsSupergiantComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AppsSupergiantComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AppsSupergiantComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

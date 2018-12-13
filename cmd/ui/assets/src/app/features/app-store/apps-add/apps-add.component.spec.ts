import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AppsAddComponent } from './apps-add.component';

describe('AppsAddComponent', () => {
  let component: AppsAddComponent;
  let fixture: ComponentFixture<AppsAddComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AppsAddComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AppsAddComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  xit('should create', () => {
    expect(component).toBeTruthy();
  });
});

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';

import { NewAppListComponent } from './new-app-list.component';

describe('NewAppListComponent', () => {
  let component: NewAppListComponent;
  let fixture: ComponentFixture<NewAppListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
    schemas: [NO_ERRORS_SCHEMA],
      declarations: [ NewAppListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NewAppListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  xit('should create', () => {
    expect(component).toBeTruthy();
  });
});

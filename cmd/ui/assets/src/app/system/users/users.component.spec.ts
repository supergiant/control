import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { Users2000Component } from './users.component';

describe('Users2000Component', () => {
  let component: Users2000Component;
  let fixture: ComponentFixture<Users2000Component>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ Users2000Component ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(Users2000Component);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

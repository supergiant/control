import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CloudAccounts2000Component } from './cloud-accounts.component';

describe('CloudAccounts2000Component', () => {
  let component: CloudAccounts2000Component;
  let fixture: ComponentFixture<CloudAccounts2000Component>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CloudAccounts2000Component ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CloudAccounts2000Component);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

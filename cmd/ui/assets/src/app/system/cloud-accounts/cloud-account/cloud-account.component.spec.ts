import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CloudAccount2000Component } from './cloud-account.component';

describe('CloudAccountComponent', () => {
  let component: CloudAccount2000Component;
  let fixture: ComponentFixture<CloudAccount2000Component>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CloudAccount2000Component ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CloudAccount2000Component);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

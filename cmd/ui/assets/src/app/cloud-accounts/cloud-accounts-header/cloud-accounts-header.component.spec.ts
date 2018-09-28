import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CloudAccountsHeaderComponent } from './cloud-accounts-header.component';

describe('CloudAccountsHeaderComponent', () => {
  let component: CloudAccountsHeaderComponent;
  let fixture: ComponentFixture<CloudAccountsHeaderComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CloudAccountsHeaderComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CloudAccountsHeaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  xit('should be created', () => {
    expect(component).toBeTruthy();
  });
});

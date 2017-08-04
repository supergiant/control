import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CloudAccountsNewSubmitComponent } from './cloud-accounts-new-submit.component';

describe('CloudAccountsNewSubmitComponent', () => {
  let component: CloudAccountsNewSubmitComponent;
  let fixture: ComponentFixture<CloudAccountsNewSubmitComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CloudAccountsNewSubmitComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CloudAccountsNewSubmitComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

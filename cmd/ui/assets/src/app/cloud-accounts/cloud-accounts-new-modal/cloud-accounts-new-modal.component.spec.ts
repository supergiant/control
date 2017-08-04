import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CloudAccountsNewModalComponent } from './cloud-accounts-new-modal.component';

describe('CloudAccountsNewModalComponent', () => {
  let component: CloudAccountsNewModalComponent;
  let fixture: ComponentFixture<CloudAccountsNewModalComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CloudAccountsNewModalComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CloudAccountsNewModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

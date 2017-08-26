import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CloudAccountDetailsComponent } from './cloud-account-details.component';

describe('CloudAccountDetailsComponent', () => {
  let component: CloudAccountDetailsComponent;
  let fixture: ComponentFixture<CloudAccountDetailsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [CloudAccountDetailsComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CloudAccountDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

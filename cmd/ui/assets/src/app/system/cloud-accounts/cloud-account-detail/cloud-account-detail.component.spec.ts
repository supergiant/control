import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CloudAccountDetailComponent } from './cloud-account-detail.component';

describe('CloudAccountDetailComponent', () => {
  let component: CloudAccountDetailComponent;
  let fixture: ComponentFixture<CloudAccountDetailComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CloudAccountDetailComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CloudAccountDetailComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

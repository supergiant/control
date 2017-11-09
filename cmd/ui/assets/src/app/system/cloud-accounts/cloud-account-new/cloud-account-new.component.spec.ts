import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CloudAccountNewComponent } from './cloud-account-new.component';

describe('CloudAccountNewComponent', () => {
  let component: CloudAccountNewComponent;
  let fixture: ComponentFixture<CloudAccountNewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CloudAccountNewComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CloudAccountNewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

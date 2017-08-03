import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CloudAccountComponent } from './cloud-account.component';

describe('CloudAccountComponent', () => {
  let component: CloudAccountComponent;
  let fixture: ComponentFixture<CloudAccountComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CloudAccountComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CloudAccountComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

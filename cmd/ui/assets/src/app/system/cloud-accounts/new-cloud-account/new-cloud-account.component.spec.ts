import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NewCloudAccountComponent } from './new-cloud-account.component';

describe('NewCloudAccountComponent', () => {
  let component: NewCloudAccountComponent;
  let fixture: ComponentFixture<NewCloudAccountComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NewCloudAccountComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NewCloudAccountComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

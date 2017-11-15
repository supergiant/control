import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { EditCloudAccountComponent } from './edit-cloud-account.component';

describe('EditCloudAccountComponent', () => {
  let component: EditCloudAccountComponent;
  let fixture: ComponentFixture<EditCloudAccountComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ EditCloudAccountComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(EditCloudAccountComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

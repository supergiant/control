import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AppsHeaderComponent } from './apps-header.component';

describe('AppsHeaderComponent', () => {
  let component: AppsHeaderComponent;
  let fixture: ComponentFixture<AppsHeaderComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AppsHeaderComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AppsHeaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

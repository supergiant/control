import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ServicesHeaderComponent } from './services-header.component';

describe('ServicesHeaderComponent', () => {
  let component: ServicesHeaderComponent;
  let fixture: ComponentFixture<ServicesHeaderComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ServicesHeaderComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ServicesHeaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

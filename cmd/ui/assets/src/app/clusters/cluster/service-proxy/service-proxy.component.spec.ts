import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ServiceProxyComponent } from './service-proxy.component';

describe('ServiceProxyComponent', () => {
  let component: ServiceProxyComponent;
  let fixture: ComponentFixture<ServiceProxyComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ServiceProxyComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ServiceProxyComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

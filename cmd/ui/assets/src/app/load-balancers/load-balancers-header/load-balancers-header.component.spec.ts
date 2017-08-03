import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { LoadBalancersHeaderComponent } from './load-balancers-header.component';

describe('LoadBalancersHeaderComponent', () => {
  let component: LoadBalancersHeaderComponent;
  let fixture: ComponentFixture<LoadBalancersHeaderComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ LoadBalancersHeaderComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LoadBalancersHeaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

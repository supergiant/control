import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { LoadBalancerDetailsComponent } from './load-balancer-details.component';

describe('LoadBalancerDetailsComponent', () => {
  let component: LoadBalancerDetailsComponent;
  let fixture: ComponentFixture<LoadBalancerDetailsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [LoadBalancerDetailsComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LoadBalancerDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

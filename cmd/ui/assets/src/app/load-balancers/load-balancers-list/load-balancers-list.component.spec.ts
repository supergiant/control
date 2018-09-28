import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { LoadBalancersListComponent } from './load-balancers-list.component';

describe('LoadBalancersListComponent', () => {
  let component: LoadBalancersListComponent;
  let fixture: ComponentFixture<LoadBalancersListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [LoadBalancersListComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LoadBalancersListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  xit('should be created', () => {
    expect(component).toBeTruthy();
  });
});

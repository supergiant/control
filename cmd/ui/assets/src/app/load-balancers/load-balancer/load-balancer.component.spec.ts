import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { LoadBalancerComponent } from './load-balancer.component';

describe('LoadBalancerComponent', () => {
  let component: LoadBalancerComponent;
  let fixture: ComponentFixture<LoadBalancerComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ LoadBalancerComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LoadBalancerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

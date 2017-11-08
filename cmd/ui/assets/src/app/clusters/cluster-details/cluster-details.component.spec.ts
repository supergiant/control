import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ClusterDetailsComponent } from './cluster-details.component';

describe('ClusterDetailsComponent', () => {
  let component: ClusterDetailsComponent;
  let fixture: ComponentFixture<ClusterDetailsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ClusterDetailsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ClusterDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ClusterListModalComponent } from './cluster-list-modal.component';

describe('ClusterListModalComponent', () => {
  let component: ClusterListModalComponent;
  let fixture: ComponentFixture<ClusterListModalComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ClusterListModalComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ClusterListModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  xit('should create', () => {
    expect(component).toBeTruthy();
  });
});

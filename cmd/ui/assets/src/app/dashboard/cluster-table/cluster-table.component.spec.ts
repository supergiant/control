import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ClusterTableComponent } from './cluster-table.component';

describe('ClusterTableComponent', () => {
  let component: ClusterTableComponent;
  let fixture: ComponentFixture<ClusterTableComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ClusterTableComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ClusterTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  xit('should create', () => {
    expect(component).toBeTruthy();
  });
});

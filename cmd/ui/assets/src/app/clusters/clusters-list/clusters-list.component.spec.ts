import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ClustersListComponent } from './clusters-list.component';

describe('ClustersListComponent', () => {
  let component: ClustersListComponent;
  let fixture: ComponentFixture<ClustersListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ClustersListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ClustersListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

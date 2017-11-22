import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ClustersContextComponent } from './clusters-context.component';

describe('ClustersContextComponent', () => {
  let component: ClustersContextComponent;
  let fixture: ComponentFixture<ClustersContextComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ClustersContextComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ClustersContextComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

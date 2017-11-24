import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ClustersTutorialComponent } from './clusters-tutorial.component';

describe('ClustersTutorialComponent', () => {
  let component: ClustersTutorialComponent;
  let fixture: ComponentFixture<ClustersTutorialComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ClustersTutorialComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ClustersTutorialComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

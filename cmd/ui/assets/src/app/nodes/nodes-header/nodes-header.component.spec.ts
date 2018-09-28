import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NodesHeaderComponent } from './nodes-header.component';

describe('NodesHeaderComponent', () => {
  let component: NodesHeaderComponent;
  let fixture: ComponentFixture<NodesHeaderComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [NodesHeaderComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NodesHeaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  xit('should be created', () => {
    expect(component).toBeTruthy();
  });
});

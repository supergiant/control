import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NodesListComponent } from './nodes-list.component';

describe('NodesListComponent', () => {
  let component: NodesListComponent;
  let fixture: ComponentFixture<NodesListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [NodesListComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NodesListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

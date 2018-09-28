import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PodsListComponent } from './pods-list.component';

describe('PodsListComponent', () => {
  let component: PodsListComponent;
  let fixture: ComponentFixture<PodsListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [PodsListComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PodsListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  xit('should be created', () => {
    expect(component).toBeTruthy();
  });
});

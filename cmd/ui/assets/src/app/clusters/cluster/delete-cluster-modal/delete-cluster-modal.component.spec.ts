import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DeleteClusterModalComponent } from './delete-cluster-modal.component';

describe('DeleteClusterModalComponent', () => {
  let component: DeleteClusterModalComponent;
  let fixture: ComponentFixture<DeleteClusterModalComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DeleteClusterModalComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DeleteClusterModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  xit('should create', () => {
    expect(component).toBeTruthy();
  });
});

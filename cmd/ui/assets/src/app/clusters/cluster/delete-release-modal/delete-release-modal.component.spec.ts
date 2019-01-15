import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DeleteReleaseModalComponent } from './delete-release-modal.component';

describe('DeleteReleaseModalComponent', () => {
  let component: DeleteReleaseModalComponent;
  let fixture: ComponentFixture<DeleteReleaseModalComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DeleteReleaseModalComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DeleteReleaseModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  xit('should create', () => {
    expect(component).toBeTruthy();
  });
});

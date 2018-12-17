import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ReleaseInfoModalComponent } from './release-info-modal.component';

describe('ReleaseInfoModalComponent', () => {
  let component: ReleaseInfoModalComponent;
  let fixture: ComponentFixture<ReleaseInfoModalComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ReleaseInfoModalComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ReleaseInfoModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  xit('should create', () => {
    expect(component).toBeTruthy();
  });
});

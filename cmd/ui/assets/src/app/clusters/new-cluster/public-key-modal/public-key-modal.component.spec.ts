import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PublicKeyModalComponent } from './public-key-modal.component';

describe('PublicKeyModalComponent', () => {
  let component: PublicKeyModalComponent;
  let fixture: ComponentFixture<PublicKeyModalComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PublicKeyModalComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PublicKeyModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  xit('should create', () => {
    expect(component).toBeTruthy();
  });
});

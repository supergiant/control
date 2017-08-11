import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DropdownModalComponent } from './dropdown-modal.component';

describe('DropdownModalComponent', () => {
  let component: DropdownModalComponent;
  let fixture: ComponentFixture<DropdownModalComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DropdownModalComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DropdownModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

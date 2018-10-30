import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MenuModalComponent } from './menu-modal.component';

describe('MenuModalComponent', () => {
  let component: MenuModalComponent;
  let fixture: ComponentFixture<MenuModalComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ MenuModalComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MenuModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { UsersHeaderComponent } from './users-header.component';

describe('UsersHeaderComponent', () => {
  let component: UsersHeaderComponent;
  let fixture: ComponentFixture<UsersHeaderComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ UsersHeaderComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(UsersHeaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

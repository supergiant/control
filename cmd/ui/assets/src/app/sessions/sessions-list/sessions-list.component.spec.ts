import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SessionsListComponent } from './sessions-list.component';

describe('SessionsListComponent', () => {
  let component: SessionsListComponent;
  let fixture: ComponentFixture<SessionsListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [SessionsListComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SessionsListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  xit('should be created', () => {
    expect(component).toBeTruthy();
  });
});

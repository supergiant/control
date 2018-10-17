import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TaskLogsComponent } from './task-logs.component';

describe('TaskLogsComponent', () => {
  let component: TaskLogsComponent;
  let fixture: ComponentFixture<TaskLogsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TaskLogsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TaskLogsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NewAppListComponent } from './new-app-list.component';

describe('NewAppListComponent', () => {
  let component: NewAppListComponent;
  let fixture: ComponentFixture<NewAppListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NewAppListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NewAppListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

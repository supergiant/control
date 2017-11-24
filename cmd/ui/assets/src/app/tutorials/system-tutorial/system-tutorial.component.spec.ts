import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SystemTutorialComponent } from './system-tutorial.component';

describe('SystemTutorialComponent', () => {
  let component: SystemTutorialComponent;
  let fixture: ComponentFixture<SystemTutorialComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SystemTutorialComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SystemTutorialComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AppsTutorialComponent } from './apps-tutorial.component';

describe('AppsTutorialComponent', () => {
  let component: AppsTutorialComponent;
  let fixture: ComponentFixture<AppsTutorialComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AppsTutorialComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AppsTutorialComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

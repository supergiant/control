import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { HelmAppComponent } from './helm-app.component';

describe('HelmAppComponent', () => {
  let component: HelmAppComponent;
  let fixture: ComponentFixture<HelmAppComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HelmAppComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HelmAppComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  xit('should be created', () => {
    expect(component).toBeTruthy();
  });
});

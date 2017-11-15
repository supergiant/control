import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HelmReposComponent } from './helm-repos.component';

describe('HelmReposComponent', () => {
  let component: HelmReposComponent;
  let fixture: ComponentFixture<HelmReposComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HelmReposComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HelmReposComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

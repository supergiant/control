import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KubectlConfigModalComponent } from './kubectl-config-modal.component';

describe('KubectlConfigModalComponent', () => {
  let component: KubectlConfigModalComponent;
  let fixture: ComponentFixture<KubectlConfigModalComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KubectlConfigModalComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KubectlConfigModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  xit('should create', () => {
    expect(component).toBeTruthy();
  });
});

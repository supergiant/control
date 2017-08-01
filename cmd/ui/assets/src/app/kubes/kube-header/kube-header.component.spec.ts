import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KubeHeaderComponent } from './kube-header.component';

describe('KubeHeaderComponent', () => {
  let component: KubeHeaderComponent;
  let fixture: ComponentFixture<KubeHeaderComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KubeHeaderComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KubeHeaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

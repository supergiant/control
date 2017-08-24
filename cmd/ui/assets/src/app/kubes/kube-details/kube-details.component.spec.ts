import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KubeDetailsComponent } from './kube-details.component';

describe('KubeDetailsComponent', () => {
  let component: KubeDetailsComponent;
  let fixture: ComponentFixture<KubeDetailsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [KubeDetailsComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KubeDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

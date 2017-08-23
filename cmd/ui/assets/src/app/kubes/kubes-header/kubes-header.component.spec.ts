import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KubesHeaderComponent } from './kubes-header.component';

describe('KubesHeaderComponent', () => {
  let component: KubesHeaderComponent;
  let fixture: ComponentFixture<KubesHeaderComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [KubesHeaderComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KubesHeaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

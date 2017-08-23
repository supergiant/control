import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KubesComponent } from './kubes.component';

describe('KubesComponent', () => {
  let component: KubesComponent;
  let fixture: ComponentFixture<KubesComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [KubesComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KubesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KubesListComponent } from './kubes-list.component';

describe('KubesListComponent', () => {
  let component: KubesListComponent;
  let fixture: ComponentFixture<KubesListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KubesListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KubesListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

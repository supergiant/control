import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { VolumesComponent } from './volumes.component';

describe('VolumesComponent', () => {
  let component: VolumesComponent;
  let fixture: ComponentFixture<VolumesComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ VolumesComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(VolumesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

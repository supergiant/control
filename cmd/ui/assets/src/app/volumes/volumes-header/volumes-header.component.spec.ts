import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { VolumesHeaderComponent } from './volumes-header.component';

describe('VolumesHeaderComponent', () => {
  let component: VolumesHeaderComponent;
  let fixture: ComponentFixture<VolumesHeaderComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ VolumesHeaderComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(VolumesHeaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  xit('should be created', () => {
    expect(component).toBeTruthy();
  });
});

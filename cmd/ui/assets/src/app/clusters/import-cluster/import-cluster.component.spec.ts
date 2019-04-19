import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ImportClusterComponent } from './import-cluster.component';

describe('ImportClusterComponent', () => {
  let component: ImportClusterComponent;
  let fixture: ComponentFixture<ImportClusterComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ImportClusterComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ImportClusterComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

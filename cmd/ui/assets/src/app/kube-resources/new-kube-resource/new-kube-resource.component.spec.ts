import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NewKubeResourceComponent } from './new-kube-resource.component';

describe('NewKubeResourceComponent', () => {
  let component: NewKubeResourceComponent;
  let fixture: ComponentFixture<NewKubeResourceComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NewKubeResourceComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NewKubeResourceComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

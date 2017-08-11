import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { KubeComponent } from './kube.component';

describe('CloudAccountComponent', () => {
  let component: KubeComponent;
  let fixture: ComponentFixture<KubeComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KubeComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KubeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { RemoveRepoDialogComponent } from './remove-repo-dialog.component';

describe('RemoveRepoDialogComponent', () => {
  let component: RemoveRepoDialogComponent;
  let fixture: ComponentFixture<RemoveRepoDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ RemoveRepoDialogComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(RemoveRepoDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

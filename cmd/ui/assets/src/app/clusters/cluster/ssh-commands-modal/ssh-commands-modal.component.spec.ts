import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SshCommandsModalComponent } from './ssh-commands-modal.component';

describe('SshCommandsModalComponent', () => {
  let component: SshCommandsModalComponent;
  let fixture: ComponentFixture<SshCommandsModalComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SshCommandsModalComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SshCommandsModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  xit('should create', () => {
    expect(component).toBeTruthy();
  });
});

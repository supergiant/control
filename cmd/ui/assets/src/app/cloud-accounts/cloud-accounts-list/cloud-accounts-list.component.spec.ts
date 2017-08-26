import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CloudAccountsListComponent } from './cloud-accounts-list.component';

describe('CloudAccountsListComponent', () => {
  let component: CloudAccountsListComponent;
  let fixture: ComponentFixture<CloudAccountsListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [CloudAccountsListComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CloudAccountsListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ListCloudAccountsComponent } from './list-cloud-accounts.component';

describe('ListCloudAccountsComponent', () => {
  let component: ListCloudAccountsComponent;
  let fixture: ComponentFixture<ListCloudAccountsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ListCloudAccountsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ListCloudAccountsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

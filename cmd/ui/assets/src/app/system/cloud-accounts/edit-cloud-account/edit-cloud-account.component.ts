import { Component, OnInit, OnDestroy, ViewEncapsulation } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription } from 'rxjs';

import { Supergiant } from '../../../shared/supergiant/supergiant.service';

@Component({
  selector: 'app-edit-cloud-account',
  templateUrl: './edit-cloud-account.component.html',
  styleUrls: ['./edit-cloud-account.component.scss'],
  encapsulation: ViewEncapsulation.None
})

export class EditCloudAccountComponent implements OnInit, OnDestroy {

  subscriptions = new Subscription();
  public accountName: string;
  public account: any;
  public saving = false;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private supergiant: Supergiant
  ) { }

  getAccount(name) {
    this.subscriptions.add(this.supergiant.CloudAccounts.get(name).subscribe(
      data => this.account = data,
      err => console.error(err)
    ))
  }

  editAccount(account) {
    if (!this.saving) {
      this.saving = true;

      this.supergiant.CloudAccounts.update(this.accountName, account).subscribe(
        res => this.router.navigate(['../../'], { relativeTo: this.route}),
        err => console.error(err)
      )
    }
  }

  cancel() {
    this.router.navigate(['../../'], { relativeTo: this.route});
  }

  ngOnInit() {
    this.accountName = this.route.snapshot.params.id;
    this.getAccount(this.accountName)
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

}

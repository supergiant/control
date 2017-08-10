import { Component, OnInit } from '@angular/core';
import { CloudAccountsService } from './cloud-accounts.service';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../shared/supergiant/supergiant.service'


@Component({
  selector: 'app-cloud-accounts',
  templateUrl: './cloud-accounts.component.html',
  styleUrls: ['./cloud-accounts.component.css']
})
export class CloudAccountsComponent implements OnInit {
  cloudAccounts: any;
  private subscription: Subscription;

  constructor(private cloudAccountsService: CloudAccountsService, private supergiant: Supergiant) { }

  ngOnInit() {
    this.getAccounts()
  }
  getAccounts() {
    this.subscription = this.supergiant.CloudAccounts.get().subscribe(
      cloudAccount=>{ this.cloudAccounts = cloudAccount.json().items},
      (err) =>{ this.cloudAccountsService.showNotification("warn", "Connection Issue.", err)})
  }



}

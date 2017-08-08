import { Component, OnInit } from '@angular/core';
import { CloudAccountsService } from './cloud-accounts.service';
import { Subscription } from 'rxjs/Subscription';

@Component({
  selector: 'app-cloud-accounts',
  templateUrl: './cloud-accounts.component.html',
  styleUrls: ['./cloud-accounts.component.css']
})
export class CloudAccountsComponent implements OnInit {
  cloudAccounts: any;
  private subscription: Subscription;

  constructor(private cloudAccountsService: CloudAccountsService) { }

  ngOnInit() {
    this.subscription = this.cloudAccountsService.getCloudAccounts().subscribe(
      cloudAccount=>{ this.cloudAccounts = cloudAccount.json().items})
  }



}

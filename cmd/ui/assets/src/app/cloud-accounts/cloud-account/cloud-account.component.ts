import { Component, OnInit, Input, NgModule } from '@angular/core';
import {BrowserModule} from '@angular/platform-browser'
import {FormsModule} from '@angular/forms'
import { CloudAccountsService } from '../cloud-accounts.service';

@Component({
  selector: '[app-cloud-account]',
  templateUrl: './cloud-account.component.html',
  styleUrls: ['./cloud-account.component.css']
})
export class CloudAccountComponent implements OnInit {
  @Input() cloudAccount: any;

  constructor(private cloudAccountsService: CloudAccountsService) { }

  ngOnInit() {
  }


}

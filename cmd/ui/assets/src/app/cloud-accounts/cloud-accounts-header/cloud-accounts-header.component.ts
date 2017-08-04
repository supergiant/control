import { Component, OnInit } from '@angular/core';
import { CloudAccountsService } from '../cloud-accounts.service';

@Component({
  selector: 'app-cloud-accounts-header',
  templateUrl: './cloud-accounts-header.component.html',
  styleUrls: ['./cloud-accounts-header.component.css']
})
export class CloudAccountsHeaderComponent implements OnInit {

  constructor(private messageService: CloudAccountsService) {}

  ngOnInit() {
  }



  sendOpen(message){
      this.messageService.openNewCloudServiceModal(message);
  }

}

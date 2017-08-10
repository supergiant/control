import { Component, OnInit } from '@angular/core';
import { CloudAccountsService } from '../cloud-accounts.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service'
import {CloudAccountsComponent} from '../cloud-accounts.component'
import {NgbModal, ModalDismissReasons} from '@ng-bootstrap/ng-bootstrap';
import { Subscription } from 'rxjs/Subscription';

@Component({
  selector: 'app-cloud-accounts-header',
  templateUrl: './cloud-accounts-header.component.html',
  styleUrls: ['./cloud-accounts-header.component.css']
})
export class CloudAccountsHeaderComponent implements OnInit {
  modalRef: any;
  private cloudAccountsSub: Subscription;
  providersObj: any;

  constructor(
    private cloudAccountsService: CloudAccountsService,
    private cloudAccountsComponant: CloudAccountsComponent,
    private supergiant: Supergiant,
    private modalService: NgbModal
    ) {}

  ngOnInit() {
  }

  ngAfterViewInit() {
    this.cloudAccountsSub = this.supergiant.CloudAccounts.schema().subscribe(
      (data) => { this.providersObj = data.json()},
      (err) => {this.cloudAccountsService.showNotification("warn", "Connection Issue.", err)});
  }

  sendOpen(message){
      this.cloudAccountsService.openNewCloudServiceModal(message);
  }
  editCloudAccount() {
    var selectedItems = this.cloudAccountsService.returnSelectedCloudAccount()

    if (selectedItems.length === 0) {
      this.cloudAccountsService.showNotification("warn", "Warning:", "No Provider Selected.")
    } else if (selectedItems.length > 1) {
      this.cloudAccountsService.showNotification("warn", "Warning:", "You cannot edit more than one provider at a time.")
    } else {
      this.providersObj.providers[selectedItems[0].provider].model = selectedItems[0]
      this.cloudAccountsService.openNewCloudServiceEditModal("Edit", selectedItems[0].provider, this.providersObj);
    }
  }
  deleteCloudAccount() {
    var selectedItems = this.cloudAccountsService.returnSelectedCloudAccount()

    for(let provider of selectedItems){
      this.supergiant.CloudAccounts.delete(provider.id).subscribe(
        (data) => {
          if (data.status >= 200 && data.status <= 299) {
            this.cloudAccountsService.showNotification("success", "Cloud Account: " + provider.name, "Deleted...")
            this.cloudAccountsComponant.getAccounts()
           }else{
            this.cloudAccountsService.showNotification("error", "Cloud Account: " + provider.name, "Error:" + data.statusText)}},
        (err) => {
          if (err) {
            this.cloudAccountsService.showNotification("error", "Cloud Account: " + provider.name, "Error:" + err)}},
      );
    }
  }

}

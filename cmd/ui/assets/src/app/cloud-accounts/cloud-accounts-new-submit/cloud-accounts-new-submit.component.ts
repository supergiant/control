import { Component, AfterViewInit, OnDestroy,ViewChild, ElementRef } from '@angular/core';
import { NgbModal, ModalDismissReasons, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { Subscription } from 'rxjs/Subscription';
import { CloudAccountsService } from '../cloud-accounts.service';
import { CloudAccountsComponent } from '../cloud-accounts.component'
import { Supergiant } from '../../shared/supergiant/supergiant.service'

@Component({
  selector: 'app-cloud-accounts-new-submit',
  templateUrl: './cloud-accounts-new-submit.component.html',
  styleUrls: ['./cloud-accounts-new-submit.component.css']
})

export class CloudAccountsNewSubmitComponent implements AfterViewInit, OnDestroy {
  private subscription: Subscription;
  private createCloudAccoutSub: Subscription;
  private cloudAccountSchema: any;
  private cloudAccountModel: any;
  private modalRef: NgbModalRef;
  private action: string;
  private providerID: any;
  @ViewChild('newCloudAccountEditModal') content: ElementRef;

  constructor(
    private modalService: NgbModal,
    private cloudAccountsService: CloudAccountsService,
    private cloudAccountsComponant: CloudAccountsComponent,
    private supergiant: Supergiant
  ) {}

  // Converts a Cloud Account Object JSON to an object, used with the "simple/advanced"
  // mode switcher.
  convertToObj(json) {
    this.cloudAccountModel = JSON.parse(json)
  }

  // Data init after load
  ngAfterViewInit() {
    // Check for messages from the new Cloud Accont dropdown, or edit button.
    this.subscription = this.cloudAccountsService.newEditModal.subscribe( message => {
      {
        // A schema object, contains:
        // .model -> Default UI settings.
        // .schema -> Rules for acceptance from the user.
        var msg = message[2]

        // The provider slug. Dynamically provided to the dropdown
        // by the supergiant schema api.
        var provider = message[1]

        // The action type. Edit (existing), Save(new)
        this.action = message[0]

        // Feed the model and schema to the UI.
        this.cloudAccountModel = msg.providers[provider].model
        this.cloudAccountSchema = msg.providers[provider].schema
        // Save the id in case of edit action type.
        this.providerID = this.cloudAccountModel.id
      };
      // open the New/Edit modal
      {this.open(this.content)};});
  }

  ngOnDestroy(){
    this.subscription.unsubscribe();
  }

  // open the New/Edit modal, save the ref so we can close it later.
  open(content) {
    this.modalRef = this.modalService.open(content);
  }

  // When the user clicks Save/Edit
  onSubmit() {

    if (this.action === "Edit") {
      console.log(this.cloudAccountModel)
    this.createCloudAccoutSub = this.supergiant.CloudAccounts.update(this.providerID, this.cloudAccountModel).subscribe(
      (data) => {
        if (data.status >= 200 && data.status <= 299) {
          this.cloudAccountsService.showNotification("success", "Cloud Account: " + this.cloudAccountModel.name, "Created...")
          this.modalRef.close()
          this.cloudAccountsComponant.getAccounts()
        }else{
          this.cloudAccountsService.showNotification("error", "Cloud Account: " + this.cloudAccountModel.name, "Error:" + data.statusText)}},
      (err) => {
        if (err) {
          this.cloudAccountsService.showNotification("error", "Cloud Account: " + this.cloudAccountModel.name, "Error:" + err)}},
    );
  } else {
    this.createCloudAccoutSub = this.supergiant.CloudAccounts.create(this.cloudAccountModel).subscribe(
      (data) => {
        if (data.status >= 200 && data.status <= 299) {
          this.cloudAccountsService.showNotification("success", "Cloud Account: " + this.cloudAccountModel.name, "Created...")
          this.modalRef.close()
          this.cloudAccountsComponant.getAccounts()
        }else{
          this.cloudAccountsService.showNotification("error", "Cloud Account: " + this.cloudAccountModel.name, "Error:" + data.statusText)}},
      (err) => {
        if (err) {
          this.cloudAccountsService.showNotification("error", "Cloud Account: " + this.cloudAccountModel.name, "Error:" + err)}},
    );
  }
  }

  private getDismissReason(reason: any): string {
    if (reason === ModalDismissReasons.ESC) {
      return 'by pressing ESC';
    } else if (reason === ModalDismissReasons.BACKDROP_CLICK) {
      return 'by clicking on a backdrop';
    } else {
      return  `with: ${reason}`;
    }
  }
}

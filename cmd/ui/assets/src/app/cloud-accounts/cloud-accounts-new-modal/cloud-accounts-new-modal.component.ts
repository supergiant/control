import { Component, AfterViewInit, OnDestroy,ViewChild, ElementRef } from '@angular/core';
import { NgbModal, ModalDismissReasons, NgbModalOptions } from '@ng-bootstrap/ng-bootstrap';
import { Subscription } from 'rxjs/Subscription';
import { CloudAccountsService } from '../cloud-accounts.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service'
import { Notifications } from '../../shared/notifications/notifications.service'


@Component({
  selector: 'app-cloud-accounts-new-modal',
  templateUrl: './cloud-accounts-new-modal.component.html',
  styleUrls: ['./cloud-accounts-new-modal.component.css']
})

export class CloudAccountsNewModalComponent implements AfterViewInit, OnDestroy{
   private subscription: Subscription;
   private cloudAccountsSub: Subscription;
   @ViewChild('newCloudAccountModal') content: ElementRef;
   modalRef: any;
   providers = [];
   providersObj: any;


   constructor(
     private modalService: NgbModal,
     private cloudAccountsService: CloudAccountsService,
     private supergiant: Supergiant,
     private notifications: Notifications,
   ) {}

   // After init, grab the schema
   ngAfterViewInit() {
     this.cloudAccountsSub = this.supergiant.CloudAccounts.schema().subscribe(
       (data) => { this.providersObj = data.json()
         // Push available providers to an array. Displayed in the dropdown.
         for(let key in this.providersObj.providers){
           this.providers.push(key)
         }
       },
       (err) => {this.notifications.display("warn", "Connection Issue.", err)});
     this.subscription = this.cloudAccountsService.newModal.subscribe( message => {if (message) {this.open(this.content)};});
   }

   ngOnDestroy(){
     this.subscription.unsubscribe();
   }

   open(content) {
     let options: NgbModalOptions = {
       size: 'sm'
     };
     this.modalRef = this.modalService.open(content, options);
   }

   // After user selects a provider from dropdown, new/edit modal
   // in Save mode it launched.
   sendOpen(message){
     this.modalRef.close();
     this.cloudAccountsService.openNewCloudServiceEditModal("Save", message, this.providersObj);
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

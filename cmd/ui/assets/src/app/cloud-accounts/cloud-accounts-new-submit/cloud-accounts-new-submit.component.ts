import { Component, OnInit, AfterViewInit, OnDestroy,ViewChild, ElementRef } from '@angular/core';
import {NgbModal, ModalDismissReasons, NgbModalRef} from '@ng-bootstrap/ng-bootstrap';
import { Subscription } from 'rxjs/Subscription';
import { CloudAccountsService } from '../cloud-accounts.service';
import dedent from "dedent";
import {FormlyFieldConfig} from 'ng-formly';
import {Validators, FormGroup} from '@angular/forms';
import { SchemaFormModule, WidgetRegistry, DefaultWidgetRegistry } from "angular2-schema-form";
import { Subject } from 'rxjs/Subject';
import {CloudAccountsComponent} from '../cloud-accounts.component'
import { Supergiant } from '../../shared/supergiant/supergiant.service'




@Component({
  selector: 'app-cloud-accounts-new-submit',
  templateUrl: './cloud-accounts-new-submit.component.html',
  styleUrls: ['./cloud-accounts-new-submit.component.css']
})



export class CloudAccountsNewSubmitComponent implements AfterViewInit, OnDestroy, OnInit {
  private subscription: Subscription;
  private createCloudAccoutSub: Subscription;
  private cloudAccoutSchemaSub: Subscription;
  private newCloudDefault: string;
  private cloudAccountSchema: any;
  private cloudAccountModel: any;
  private modalRef: NgbModalRef;
  private action: string;
  private providerID: any;
  cloudAccountJSON = new Subject<any>();
  @ViewChild('newCloudAccountEditModal') content: ElementRef;

  constructor(
    private modalService: NgbModal,
    private cloudAccountsService: CloudAccountsService,
    private cloudAccountsComponant: CloudAccountsComponent,
    private supergiant: Supergiant
  ) {}

  ngOnInit() {

  }

  convertToObj(json) {
    this.cloudAccountModel = JSON.parse(json)
  }

  ngAfterViewInit() {
    this.subscription = this.cloudAccountsService.newEditModal.subscribe( message => {
      {
        var msg = message[2]
        var provider = message[1]
        this.action = message[0]
        this.cloudAccountModel = msg.providers[provider].model
        this.cloudAccountSchema = msg.providers[provider].schema
        this.providerID = this.cloudAccountModel.id
      };
      {this.open(this.content)};});
  }

  ngOnDestroy(){
    this.subscription.unsubscribe();
  }

  open(content) {
    this.modalRef = this.modalService.open(content);
  }

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

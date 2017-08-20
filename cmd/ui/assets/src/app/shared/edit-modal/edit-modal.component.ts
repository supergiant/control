import { Component, OnInit, AfterViewInit, OnDestroy, ViewChild, ElementRef } from '@angular/core';
import { NgbModal, ModalDismissReasons, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { Subscription } from 'rxjs/Subscription';
import { EditModalService } from './edit-modal.service';
import { Notifications } from '../../shared/notifications/notifications.service';

@Component({
  selector: 'app-edit-modal',
  templateUrl: './edit-modal.component.html',
  styleUrls: ['./edit-modal.component.css']
})
export class EditModalComponent implements OnInit, AfterViewInit, OnDestroy {
  private modalRef: NgbModalRef;
  private subscription: Subscription;
  private schema: any;
  private model: any;
  private item: any;
  private schemaBlob: any;
  private action: string;
  private title: string;
  private textStatus = 'form-control';
  private badString: string;
  private isDisabled: boolean;
  @ViewChild('editModal') content: ElementRef;


  constructor(
    private modalService: NgbModal,
    private editModalService: EditModalService,
    private notifications: Notifications,
  ) { }


  ngOnInit() {
  }

  // Data init after load
  ngAfterViewInit() {
    // Check for messages from the new Cloud Accont dropdown, or edit button.
    this.subscription = this.editModalService.newModal.subscribe(message => {
      {
        // A schema object, contains:
        // .model -> Default UI settings.
        // .schema -> Rules for acceptance from the user.
        this.schemaBlob = message[2];

        // The item slug.
        this.item = message[1];

        // The action type. Edit (existing), Save(new)
        this.action = message[0];

        // Feed the model and schema to the UI.
        this.model = this.schemaBlob[this.item].model;
        this.schema = this.schemaBlob[this.item].schema;
        this.isDisabled = false;
        this.textStatus = 'form-control goodTextarea';
      }
      // open the New/Edit modal
      { this.open(this.content); }
    });
  }

  convertToObj(json) {
    try {
      JSON.parse(json);
    } catch (e) {
      this.textStatus = 'form-control badTextarea';
      this.badString = 'Not Valid JSON: May god have mercy on your soul...';
      this.isDisabled = true;
      return;
    }
    this.textStatus = 'form-control goodTextarea';
    this.badString = 'Nice Looking JSON Dude!';
    this.isDisabled = false;
    this.model = JSON.parse(json);
  }

  open(content) {
    const options: NgbModalOptions = {
      size: 'lg'
    };
    this.modalRef = this.modalService.open(content, options);
    // If user cancels or closes the window.. we need to answer the promise.
    this.modalRef.result.then(
      (window) => { this.editModalService.editModalResponse.next('closed'); },
      (err) => { this.editModalService.editModalResponse.next('closed'); },
    );
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }

  onSubmit() {
    this.modalRef.close();
    this.editModalService.editModalResponse.next([this.action, this.item, this.model]);
  }

}

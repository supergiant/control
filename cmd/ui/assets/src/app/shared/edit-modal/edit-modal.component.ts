import { Component, OnInit, AfterViewInit, OnDestroy,ViewChild, ElementRef } from '@angular/core';
import { NgbModal, ModalDismissReasons, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { Subscription } from 'rxjs/Subscription';
import { EditModalService } from './edit-modal.service';
import { Notifications } from '../../shared/notifications/notifications.service'

@Component({
  selector: 'app-edit-modal',
  templateUrl: './edit-modal.component.html',
  styleUrls: ['./edit-modal.component.css']
})
export class EditModalComponent implements OnInit {
  private modalRef: NgbModalRef;
  private subscription: Subscription;
  private schema: any;
  private model: any;
  private item: any;
  private schemaBlob: any;
  private action: string;
  private title: string;
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
    this.subscription = this.editModalService.newModal.subscribe( message => {
      {
        // A schema object, contains:
        // .model -> Default UI settings.
        // .schema -> Rules for acceptance from the user.
        this.schemaBlob = message[2]

        // The item slug.
        this.item = message[1]

        // The action type. Edit (existing), Save(new)
        this.action = message[0]

        // Feed the model and schema to the UI.
        this.model = this.schemaBlob.providers[this.item].model
        this.schema = this.schemaBlob.providers[this.item].schema
      };
      // open the New/Edit modal
      {this.open(this.content)};});
  }

  convertToObj(json) {
    this.model = JSON.parse(json)
  }

  open(content) {
    let options: NgbModalOptions = {
      size: 'lg'
    };
    this.modalRef = this.modalService.open(content, options);
    this.modalRef.result.then((result) => {
    });
  }

  ngOnDestroy(){
    this.subscription.unsubscribe();
  }

  onSubmit() {
    this.modalRef.close();
    this.editModalService.editModalResponse.next([this.action, this.item, this.model])
  }
}

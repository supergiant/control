import { Component, OnInit, AfterViewInit, OnDestroy,ViewChild, ElementRef } from '@angular/core';
import { NgbModal, ModalDismissReasons, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { Subscription } from 'rxjs/Subscription';
import { DropdownModalService } from './dropdown-modal.service';
import { Notifications } from '../../shared/notifications/notifications.service'

@Component({
  selector: 'app-dropdown-modal',
  templateUrl: './dropdown-modal.component.html',
  styleUrls: ['./dropdown-modal.component.css']
})
export class DropdownModalComponent implements OnInit {
  private title: string;
  private type: string;
  private options = new Array();
  private modalRef: NgbModalRef;
  private subscription: Subscription;
  @ViewChild('dropdownModal') content: ElementRef;


  constructor(
    private modalService: NgbModal,
    private dropdownModalService: DropdownModalService,
    private notifications: Notifications,
  ) { }


  ngOnInit() {
  }

  ngAfterViewInit() {
    this.subscription = this.dropdownModalService.newModal.subscribe(
      message => {if (message) {
        this.title = message[0]
        this.type = message[1]
        this.options = message[2]

        this.open(this.content)
      };});
  }

  open(content) {
    let options: NgbModalOptions = {
      size: 'sm'
    };
    console.log(this.subscription)
    this.modalRef = this.modalService.open(content, options)
    this.modalRef.result.then((result) => {
      this.modalRef = void 0
    });
  }

  ngOnDestroy(){
    this.subscription.unsubscribe();
  }

  sendChoice(choice){
    this.modalRef.close();
    this.modalRef = void 0
    this.dropdownModalService.dropdownModalResponse.next(choice)
  }
}

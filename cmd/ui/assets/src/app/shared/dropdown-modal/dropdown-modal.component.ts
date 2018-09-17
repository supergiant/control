import { Component, AfterViewInit, OnDestroy, ViewChild, ElementRef } from '@angular/core';
import { NgbModal, ModalDismissReasons, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { Subscription } from 'rxjs/Subscription';
import { DropdownModalService } from './dropdown-modal.service';
import { Notifications } from '../../shared/notifications/notifications.service';

@Component({
  selector: 'app-dropdown-modal',
  templateUrl: './dropdown-modal.component.html',
  styleUrls: ['./dropdown-modal.component.scss']
})
export class DropdownModalComponent implements AfterViewInit, OnDestroy {
  private title: string;
  private type: string;
  private options = new Array();
  private modalRef: NgbModalRef;
  subscriptions = new Subscription();
  @ViewChild('dropdownModal') content: ElementRef;


  constructor(
    private modalService: NgbModal,
    private dropdownModalService: DropdownModalService,
    private notifications: Notifications,
  ) { }

  ngAfterViewInit() {
    this.subscriptions.add(this.dropdownModalService.newModal.subscribe(
      message => {
        if (message) {
          this.title = message[0];
          this.type = message[1];
          this.options = message[2];

          this.open(this.content);
        }
      }));
  }

  open(content) {
    const options: NgbModalOptions = {
      size: 'sm'
    };
    this.modalRef = this.modalService.open(content, options);
    // If user cancels or closes the window.. we need to answer the promise.
    this.modalRef.result.then(
      (window) => { this.dropdownModalService.dropdownModalResponse.next('closed'); },
      (err) => { this.dropdownModalService.dropdownModalResponse.next('closed'); },
    );
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  sendChoice(choice) {
    this.modalRef.close();
    this.dropdownModalService.dropdownModalResponse.next(choice);
  }
}

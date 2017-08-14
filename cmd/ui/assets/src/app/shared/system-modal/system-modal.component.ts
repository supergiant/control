import { Component, OnDestroy,ViewChild, ElementRef } from '@angular/core';
import { NgbModal, NgbModalOptions } from '@ng-bootstrap/ng-bootstrap';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../../shared/supergiant/supergiant.service'
import { SystemModalService } from './system-modal.service'
import { Notifications } from '../notifications/notifications.service'
import { Observable } from 'rxjs/Rx';
import {Http} from '@angular/http';

@Component({
  selector: 'app-system-modal',
  templateUrl: './system-modal.component.html',
  styleUrls: ['./system-modal.component.css']
})
export class SystemModalComponent {
  private subscription: Subscription;
  private logstream: Subscription;
  private logData: any;
  @ViewChild('systemModal') content: ElementRef;


  constructor(
    private modalService: NgbModal,
    private supergiant: Supergiant,
    private systemModalService: SystemModalService,
    private notifications: Notifications,
    http: Http,
  ) {}

  // After init, grab the subscription.
  ngAfterViewInit() {
    this.subscription = this.systemModalService.newModal.subscribe(
      message => {if (message) {this.open(this.content)};});
  }

  ngOnDestroy(){
    this.logstream.unsubscribe();
    this.subscription.unsubscribe();
  }

  open(content) {
    this.logstream = Observable.timer(0, 1000)
    .switchMap(() => this.supergiant.Logs.get()).subscribe(
      (data) => {
        this.logData = data.text()
        this.logData = this.logData.replace(/[\x00-\x7F]\[\d+mINFO[\x00-\x7F]\[0m/g, "<span class='text-info'>INFO</span> ")
        this.logData = this.logData.replace(/[\x00-\x7F]\[\d+mWARN[\x00-\x7F]\[0m/g, "<span class='text-warning'>WARN</span> ")
        this.logData = this.logData.replace(/[\x00-\x7F]\[\d+mERRO[\x00-\x7F]\[0m/g, "<span class='text-danger'>ERRO</span> ")
        this.logData = this.logData.replace(/[\x00-\x7F]\[\d+mDEBU[\x00-\x7F]\[0m/g, "<span class='text-muted'>DEBU</span> ")
      },
      );

    let options: NgbModalOptions = {
      size: 'lg'
    };
    this.modalService.open(content, options);
  }
}

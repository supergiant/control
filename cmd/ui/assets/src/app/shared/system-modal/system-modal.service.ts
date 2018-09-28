import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';

@Injectable()
export class SystemModalService {
  newModal = new Subject<any>();
  notifications = new Array();

  constructor() { }

  openSystemModal(message) {
    this.newModal.next(message);
  }

  recordNotification(notification) {
    this.notifications.push(notification);
  }
}

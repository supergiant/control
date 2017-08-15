import { Injectable } from '@angular/core';
import { Subject } from 'rxjs/Subject';

@Injectable()
export class SystemModalService {
    newModal = new Subject<any>();
    notifications = new Array();

    constructor() {}

    openSystemModal(message){
      this.newModal.next(message);
    }

    recordNotification(notification){
      this.notifications.push(notification)
    }
}

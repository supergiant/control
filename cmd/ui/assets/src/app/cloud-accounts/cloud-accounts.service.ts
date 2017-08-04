import { Injectable } from '@angular/core';
import { Subject } from 'rxjs/Subject';

@Injectable()
export class CloudAccountsService {
    newModal = new Subject<any>();

    constructor() { }

    openNewCloudServiceModal(message){
      console.log("sending " + message);
      this.newModal.next(message);
    }
}

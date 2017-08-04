import { Injectable } from '@angular/core';
import { Subject } from 'rxjs/Subject';

@Injectable()
export class CloudAccountsService {
    newModal = new Subject<any>();
    newEditModal = new Subject<any>();

    constructor() { }

    openNewCloudServiceModal(message){
      this.newModal.next(message);
    }
    openNewCloudServiceEditModal(message){
      console.log("sending " + message);
      this.newEditModal.next(message);
    }
}

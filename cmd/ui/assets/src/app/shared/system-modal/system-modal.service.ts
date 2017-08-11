import { Injectable } from '@angular/core';
import { Subject } from 'rxjs/Subject';

@Injectable()
export class SystemModalService {
    newModal = new Subject<any>();

    constructor() {}

    openSystemModal(message){
      this.newModal.next(message);
    }
}

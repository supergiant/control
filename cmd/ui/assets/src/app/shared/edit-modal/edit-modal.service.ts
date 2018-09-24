import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';

@Injectable()
export class EditModalService {
  newModal = new Subject<any>();
  editModalResponse = new Subject<any>();

  constructor() { }

  open(type, message, object) {
    this.newModal.next([type, message, object]);
  }
}

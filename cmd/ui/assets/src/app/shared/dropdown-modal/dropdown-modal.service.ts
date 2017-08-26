import { Injectable } from '@angular/core';
import { Subject } from 'rxjs/Subject';

@Injectable()
export class DropdownModalService {
  newModal = new Subject<any>();
  dropdownModalResponse = new Subject<any>();

  constructor() { }

  open(title, type, options) {
    this.newModal.next([title, type, options]);
  }
}

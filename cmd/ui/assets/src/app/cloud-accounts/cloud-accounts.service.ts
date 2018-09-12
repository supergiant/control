import { Injectable } from '@angular/core';
import { Subject ,  Observable } from 'rxjs';

@Injectable()
export class CloudAccountsService {
  newModal = new Subject<any>();
  newEditModal = new Subject<any>();
  cloudAccounts: any;
  selectedItems = new Array();
  searchString: string;

  constructor() { }

  // return all selected cloud accounts
  returnSelected() {
    return this.selectedItems;
  }

  isChecked(item) {
    for (const obj of this.selectedItems) {
      if (item.id === obj.id) { return true; }
    }
    return false;
  }

  resetSelected() {
    this.selectedItems = [];
  }
  // Record/Delete a cloud account selection from the "selected items" array.
  selectItem(item, event) {
    if (event) {
      this.selectedItems.push(item);
    } else {
      for (const obj of this.selectedItems) {
        if (item.id === obj.id) {
          this.selectedItems.splice(
            this.selectedItems.indexOf(obj), 1);
        }
      }
    }
  }
}

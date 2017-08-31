import { Injectable } from '@angular/core';
import { Subject } from 'rxjs/Subject';

@Injectable()
export class RepoModalService {
  newModal = new Subject<any>();
  repoModalResponse = new Subject<any>();
  selectedItems = new Array();

  constructor() { }

  openRepoModal(message) {
    this.newModal.next(message);
  }

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

  // Record/Delete a selection from the "selected items" array.
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

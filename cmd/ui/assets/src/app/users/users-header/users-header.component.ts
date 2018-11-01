import { Component, OnDestroy, AfterViewInit } from '@angular/core';
import { UsersService } from '../users.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Subscription } from 'rxjs';
import { Notifications } from '../../shared/notifications/notifications.service';
import { SystemModalService } from '../../shared/system-modal/system-modal.service';
import { DropdownModalService } from '../../shared/dropdown-modal/dropdown-modal.service';
import { EditModalService } from '../../shared/edit-modal/edit-modal.service';
import { UsersModel } from '../users.model';
import { LoginComponent } from '../../login/login.component';



@Component({
  selector: 'app-users-header',
  templateUrl: './users-header.component.html',
  styleUrls: ['./users-header.component.scss']
})
export class UsersHeaderComponent implements OnDestroy, AfterViewInit {
  providersObj: any;
  subscriptions = new Subscription();
  editID: number;
  searchString = '';
  constructor(
    private usersService: UsersService,
    private supergiant: Supergiant,
    private notifications: Notifications,
    private systemModalService: SystemModalService,
    private dropdownModalService: DropdownModalService,
    private editModalService: EditModalService,
    public loginComponent: LoginComponent
  ) { }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  setSearch(value) {
    this.usersService.searchString = value;
  }

  ngAfterViewInit() {
    this.subscriptions.add(this.editModalService.editModalResponse.subscribe(
      (userInput) => {
        if (userInput !== 'close') {
          const action = userInput[0];
          const providerID = userInput[1];
          const model = userInput[2];

          if (action === 'Save') {
            this.subscriptions.add(this.supergiant.Users.create(model).subscribe(
              (data) => {
                this.success(model);
                this.usersService.resetSelected();
              },
              (err) => { this.error(model, err); }
            ));
          } else if (action === 'Edit') {
            this.subscriptions.add(this.supergiant.Users.update(this.editID, model).subscribe(
              (data) => {
                this.success(model);
                this.usersService.resetSelected();
              },
              (err) => { this.error(model, err); }
            ));
          }
        }
      }
    ));
  }

  success(model) {
    this.notifications.display(
      'success',
      'User: ' + model.username,
      'Created...',
    );
  }

  error(model, data) {
    this.notifications.display(
      'error',
      'User: ' + model.username,
      'Error:' + data.statusText);
  }

  // If new button if hit, the New dropdown is triggered.
  newUser(message) {
    const userModel = new UsersModel;
    this.editModalService.open('Save', 'user', userModel);
  }

  openSystemModal(message) {
    this.systemModalService.openSystemModal(message);
  }
  // If the edit button is hit, the Edit modal is opened.
  editUser() {
    const userModel = new UsersModel;
    const selectedItems = this.usersService.returnSelected();

    if (selectedItems.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No User Selected.');
    } else if (selectedItems.length > 1) {
      this.notifications.display('warn', 'Warning:', 'You cannot edit more than one User at a time.');
    } else {
      this.editID = selectedItems[0].id;
      userModel.user['model'] = selectedItems[0];
      this.editModalService.open('Edit', 'user', userModel);
    }
  }

  generateApiToken() {
    const selectedItems = this.usersService.returnSelected();

    if (selectedItems.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No User Selected.');
    } else {
      for (const user of selectedItems) {
        this.subscriptions.add(this.supergiant.Users.generateToken(user.id).subscribe(
          (data) => {
            this.notifications.display('success', 'User: ' + user.username, 'API Key Updated...');
            this.usersService.resetSelected();
          },
          (err) => {
            this.notifications.display('error', 'User: ' + user.username, 'Error:' + err);
          },
        ));
      }
    }
  }
  // If the delete button is hit, the seleted accounts are deleted.
  deleteUser() {
    const selectedItems = this.usersService.returnSelected();
    if (selectedItems.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No User Selected.');
    } else {
      for (const user of selectedItems) {
        this.subscriptions.add(this.supergiant.Users.delete(user.id).subscribe(
          (data) => {
            this.notifications.display('success', 'User: ' + user.username, 'Deleted...');
            this.usersService.resetSelected();
          },
          (err) => {
            this.notifications.display('error', 'User: ' + user.username, 'Error:' + err);
          },
        ));
      }
    }
  }
}

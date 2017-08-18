import { Component } from '@angular/core';
import { UsersService } from '../users.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service'
import { UsersComponent } from '../users.component'
import { Subscription } from 'rxjs/Subscription';
import { Notifications } from '../../shared/notifications/notifications.service'
import { SystemModalService } from '../../shared/system-modal/system-modal.service'
import { DropdownModalService } from '../../shared/dropdown-modal/dropdown-modal.service'
import { EditModalService } from '../../shared/edit-modal/edit-modal.service'
import { UsersModel } from '../users.model';



@Component({
  selector: 'app-users-header',
  templateUrl: './users-header.component.html',
  styleUrls: ['./users-header.component.css']
})
export class UsersHeaderComponent {
  providersObj: any;
  subscriptions = [];
  editID: number;
  constructor(
    private usersService: UsersService,
    private usersComponent: UsersComponent,
    private supergiant: Supergiant,
    private notifications: Notifications,
    private systemModalService: SystemModalService,
    private dropdownModalService: DropdownModalService,
    private editModalService: EditModalService,
  ) { }

  ngOnDestroy() {
    for (let subscription of this.subscriptions) {
      subscription.unsubscribe();
    }
  }

  ngAfterViewInit() {
    this.subscriptions["edit"] = this.editModalService.editModalResponse.subscribe(
      (userInput) => {
        if (userInput != "close") {
          var action = userInput[0]
          var providerID = 1
          var model = userInput[2]

          if (action === "Save") {
            this.supergiant.Users.create(model).subscribe(
              (data) => {
                this.success(model)
                this.usersComponent.getUsers()
              },
              (err) => { this.error(model, err) }
            );
          } else if (action === "Edit") {
            this.supergiant.Users.update(this.editID, model).subscribe(
              (data) => {
                this.success(model)
                this.usersService.resetSelected()
                this.usersComponent.getUsers()
              },
              (err) => { this.error(model, err) }
            );
          }
        }
      }
    );
  }

  success(model) {
    this.notifications.display(
      "success",
      "User: " + model.name,
      "Created...",
    )
  }

  error(model, data) {
    this.notifications.display(
      "error",
      "User: " + model.name,
      "Error:" + data.statusText)
  }

  // If new button if hit, the New dropdown is triggered.
  newUser(message) {
    let userModel = new UsersModel
    this.editModalService.open("Save", 'user', userModel)
  }

  openSystemModal(message) {
    this.systemModalService.openSystemModal(message);
  }
  // If the edit button is hit, the Edit modal is opened.
  editUser() {
    let userModel = new UsersModel
    var selectedItems = this.usersService.returnSelected()

    if (selectedItems.length === 0) {
      this.notifications.display("warn", "Warning:", "No User Selected.")
    } else if (selectedItems.length > 1) {
      this.notifications.display("warn", "Warning:", "You cannot edit more than one User at a time.")
    } else {
      this.editID = selectedItems[0].id
      userModel.user["model"] = selectedItems[0]
      this.editModalService.open("Edit", 'user', userModel);
    }
  }

  generateApiToken() {
    var selectedItems = this.usersService.returnSelected()

    if (selectedItems.length === 0) {
      this.notifications.display("warn", "Warning:", "No User Selected.")
    } else {
      for (let user of selectedItems) {
        this.subscriptions[user.id] = this.supergiant.Users.generateToken(user.id).subscribe(
          (data) => {
            this.notifications.display("success", "User: " + user.username, "API Key Updated...")
            this.usersService.resetSelected()
            this.usersComponent.getUsers()
          },
          (err) => {
            this.notifications.display("error", "User: " + user.username, "Error:" + err)
          },
        )
      }
    }
  }
  // If the delete button is hit, the seleted accounts are deleted.
  deleteUser() {
    var selectedItems = this.usersService.returnSelected()
    if (selectedItems.length === 0) {
      this.notifications.display("warn", "Warning:", "No User Selected.")
    } else {
      for (let user of selectedItems) {
        this.supergiant.Users.delete(user.id).subscribe(
          (data) => {
            this.notifications.display("success", "User: " + user.username, "Deleted...")
            this.usersService.resetSelected()
            this.usersComponent.getUsers()
          },
          (err) => {
            this.notifications.display("error", "User: " + user.username, "Error:" + err)
          },
        );
      }
    }
  }
}

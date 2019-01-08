import { Injectable } from '@angular/core';
import { NotificationsService } from 'angular2-notifications';


@Injectable()
export class Notifications {
  constructor(
    private _service: NotificationsService,
  ) { }

  // Notification Shortcut
  display(kind, header, body) {
    switch (kind) {
      case 'success': {
        this._service.success(header, body, {});
        break;
      }
      case 'error': {
        this._service.error(header, body, {});
        break;
      }
      case 'warn': {
        this._service.warn(header, body, {});
        break;
      }
    }
  }
}

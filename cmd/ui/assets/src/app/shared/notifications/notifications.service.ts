import { Injectable } from '@angular/core';
import { NotificationsService } from 'angular2-notifications/dist';
import { SystemModalService } from '../system-modal/system-modal.service';


@Injectable()
export class Notifications {
  constructor(
    private _service: NotificationsService,
    private systemModalService: SystemModalService,
  ) { }

  // Notification Shortcut
  display(kind, header, body) {
    switch (kind) {
      case 'success': {
        this._service.success(header, body, {});
        this.systemModalService.recordNotification([kind, body]);
        break;
      }
      case 'error': {
        this._service.error(header, body, {});
        this.systemModalService.recordNotification([kind, body]);
        break;
      }
      case 'warn': {
        this._service.warn(header, body, {});
        this.systemModalService.recordNotification([kind, body]);
        break;
      }
    }
  }
}

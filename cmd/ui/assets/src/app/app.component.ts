import { Component } from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'app';

  public options = {
    position: ['top', 'left'],
    timeOut: 2000,
    lastOnBottom: true,
  };
}

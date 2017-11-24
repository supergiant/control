import { Component } from '@angular/core';
import {ActivatedRoute} from '@angular/router';




@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'app';
  public location = '' ;
  public options = {
    position: ['top', 'left'],
    timeOut: 2000,
    lastOnBottom: true,
  };
  constructor(private  _router: ActivatedRoute) {
    console.log(_router);
    this.location = _router.snapshot.url.join('');
  }
}

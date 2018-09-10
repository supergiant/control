import { Component, ViewChild } from '@angular/core';
import { Router } from '@angular/router';


@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {

  @ViewChild('body') body;

  title = 'app';
  // TODO: generate this dynamically based on img/backgrounds count
  public bg_count = 12;
  public location = '';
  public options = {
    position: ['top', 'left'],
    timeOut: 2000,
    lastOnBottom: true,
  };

  constructor( private router: Router) {
  }

  ngAfterViewInit() {
    // const i = Math.floor(Math.random() * this.bg_count) + 1;
    // const body = this.body.nativeElement;
    // // TODO: figure out resizing imgs
    // body.style.background = "url(assets/img/backgrounds/bg-" + i.toString() + ".jpg) center top / 1366px 768px no-repeat";
  }

  get showHeader() {
    return !Boolean(this.router.url === '/');
  }
}

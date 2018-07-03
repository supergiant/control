import { Component } from '@angular/core';
import { ActivatedRoute, Router, Event, NavigationEnd } from '@angular/router';




@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'app';
  public showHeader = true;
  public location = '';
  public options = {
    position: ['top', 'left'],
    timeOut: 2000,
    lastOnBottom: true,
  };

  constructor(private _route: ActivatedRoute, router: Router) {
    this.location = _route.snapshot.url.join('');
    router.events.subscribe((e: Event) => { // TODO: think about implementing route guard instead
      if (e instanceof NavigationEnd) {
        if (e.urlAfterRedirects === "/" || e.urlAfterRedirects === "/#log-out") {
          this.showHeader = false;
        } else {
          this.showHeader = true;
        }
      }
    })
  }
}

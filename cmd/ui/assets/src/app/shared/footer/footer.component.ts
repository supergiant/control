import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { of } from 'rxjs';
import { catchError } from 'rxjs/operators';

@Component({
  selector: 'app-footer',
  templateUrl: './footer.component.html',
  styleUrls: ['./footer.component.scss']
})
export class FooterComponent implements OnInit {

  version$: any;

  constructor(private http: HttpClient) {}

  ngOnInit() {
    this.version$ = this.http.get('/version', { responseType: 'text' })
      .pipe(
        catchError(err => {
          console.error(err)
          return of('');
        })
      )
  }
}

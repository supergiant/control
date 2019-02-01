import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-footer',
  templateUrl: './footer.component.html',
  styleUrls: ['./footer.component.scss']
})
export class FooterComponent implements OnInit {

  version: any;

  constructor(private http: HttpClient) {}

  ngOnInit() {
    this.http.get('/version', { responseType: 'text' }).subscribe(
      res => this.version = res,
      err => console.error(err)
    )
  }
}

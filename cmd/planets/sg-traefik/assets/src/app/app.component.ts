import { Component, OnInit } from '@angular/core';
import { Http, Headers } from '@angular/http';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  title = 'app';
  traefikURL = 'http://aed2af6bf904511e7a1470293ee4ef66-1761795073.us-east-1.elb.amazonaws.com';
  healthAPI = '/health';
  providerAPI = '/api';
  public pieChartLabels: string[];
  public pieChartData: number[];
  public pieChartType = 'doughnut';
  isDataAvailable = false;
  healthData: any;

  constructor(
    private http: Http,
  ) { }

  ngOnInit() {
    // const headers = new Headers();
    // headers.append('Access-Control-Allow-Origin', '*');
    this.http.get(this.traefikURL + this.healthAPI).subscribe(
      (health) => {
        this.healthData = health.json();
        const metrics = health.json().total_status_code_count;
        if (metrics) {
          this.isDataAvailable = true;
          this.pieChartLabels = Object.keys(metrics);
          this.pieChartData = Object.values(metrics);
        }

      },
      (err) => { console.log(err); }
    );
  }

}

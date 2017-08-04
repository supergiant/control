import { Component, AfterViewInit, OnDestroy,ViewChild, ElementRef } from '@angular/core';
import {NgbModal, ModalDismissReasons} from '@ng-bootstrap/ng-bootstrap';
import { Subscription } from 'rxjs/Subscription';
import { CloudAccountsService } from '../cloud-accounts.service';
import dedent from "dedent";

@Component({
  selector: 'app-cloud-accounts-new-submit',
  templateUrl: './cloud-accounts-new-submit.component.html',
  styleUrls: ['./cloud-accounts-new-submit.component.css']
})
export class CloudAccountsNewSubmitComponent implements AfterViewInit, OnDestroy {
  private subscription: Subscription;
  private newCloudDefault: string;
  @ViewChild('newCloudAccountEditModal') content: ElementRef;



  constructor(private modalService: NgbModal, private cloudAccountsService: CloudAccountsService) {}

  ngAfterViewInit() {
    this.subscription = this.cloudAccountsService.newEditModal.subscribe( message => {
      {
        switch(message) {
          case "aws": {
          this.newCloudDefault = dedent(`
            {
              "credentials": {
                "access_key": "",
                "secret_key": ""
              },
              "name": "",
              "provider": "aws"
            }
                                 `);
          break;
          }
          case "do": {
          this.newCloudDefault = dedent(`
            {
              "credentials": {
                "token": ""
              },
              "name": "",
              "provider": "digitalocean"
            }
                                 `);
          break;
          }
          case "gce": {
          this.newCloudDefault =dedent(`
            {
              "credentials": {
                "auth_provider_x509_cert_url": "",
                "auth_uri": "",
                "client_email": "",
                "client_id": "",
                "client_x509_cert_url": "",
                "private_key": "",
                "private_key_id": "",
                "project_id": "",
                "token_uri": "",
                "type": ""
              },
              "name": "",
              "provider": "gce"
            }
                                 `);
          break;
          }
          case "os": {
          this.newCloudDefault =dedent(`
            {
              "credentials": {
                "domain_id": "",
                "domain_name": "",
                "identity_endpoint": "",
                "password": "",
                "tenant_id": "",
                "username": ""
              },
              "name": "",
              "provider": "openstack"
            }
                                 `);
          break;
          }
          case "pkt": {
          this.newCloudDefault =dedent(`
            {
              "credentials": {
                "api_token": ""
              },
              "name": "",
              "provider": "packet"
            }
                                 `);
          break;
          }
          default: {
          this.newCloudDefault = "No Data"
          break;
          }
        }
      };
      {this.open(this.content)};});
  }

  ngOnDestroy(){
    this.subscription.unsubscribe();
  }

  open(content) {
    this.modalService.open(content);
  }

  sendOpen(message){
      this.cloudAccountsService.openNewCloudServiceEditModal(message);
  }

  private getDismissReason(reason: any): string {
    if (reason === ModalDismissReasons.ESC) {
      return 'by pressing ESC';
    } else if (reason === ModalDismissReasons.BACKDROP_CLICK) {
      return 'by clicking on a backdrop';
    } else {
      return  `with: ${reason}`;
    }
  }

}

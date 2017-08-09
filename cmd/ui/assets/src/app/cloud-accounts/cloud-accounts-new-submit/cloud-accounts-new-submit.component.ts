import { Component, OnInit, AfterViewInit, OnDestroy,ViewChild, ElementRef } from '@angular/core';
import {NgbModal, ModalDismissReasons, NgbModalRef} from '@ng-bootstrap/ng-bootstrap';
import { Subscription } from 'rxjs/Subscription';
import { CloudAccountsService } from '../cloud-accounts.service';
import dedent from "dedent";
import {FormlyFieldConfig} from 'ng-formly';
import {Validators, FormGroup} from '@angular/forms';
import { SchemaFormModule, WidgetRegistry, DefaultWidgetRegistry } from "angular2-schema-form";
import { Subject } from 'rxjs/Subject';




@Component({
  selector: 'app-cloud-accounts-new-submit',
  templateUrl: './cloud-accounts-new-submit.component.html',
  styleUrls: ['./cloud-accounts-new-submit.component.css']
})



export class CloudAccountsNewSubmitComponent implements AfterViewInit, OnDestroy, OnInit {
  private subscription: Subscription;
  private newCloudDefault: string;
  private cloudAccountSchema: any;
  private cloudAccountModel: any;
  private modalRef: NgbModalRef;
  cloudAccountJSON = new Subject<any>();
  @ViewChild('newCloudAccountEditModal') content: ElementRef;

  constructor(private modalService: NgbModal, private cloudAccountsService: CloudAccountsService) {}

  ngOnInit() {

  }

  convertToObj(json) {
    this.cloudAccountModel = JSON.parse(json)
  }

  ngAfterViewInit() {
    this.subscription = this.cloudAccountsService.newEditModal.subscribe( message => {
      {
        switch(message) {
          case "aws": {
          this.cloudAccountModel = {
              "credentials": {
                "access_key": "",
                "secret_key": ""
              },
              "name": "",
              "provider": "aws"
            }

            this.cloudAccountSchema = {
              "properties": {
                "credentials": {
                  "type": "object",
                  "properties": {
                    "access_key": {
                      "type": "string",
                      "description": "Access Key"
                    },
                    "secret_key": {
                      "type": "string",
                      "description": "Secret Access Key"
                    }
                  }
                },
                "name": {
                  "type": "string",
                  "description": "Provider Name"
                },
                "provider": {
                  "type": "string",
                  "default": "aws",
                  "description": "Provider",
                  "widget": "hidden"
                }
              }
            }
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
    this.modalRef = this.modalService.open(content);
  }

  onSubmit() {
    var err = this.cloudAccountsService.createCloudAccount(JSON.stringify(this.cloudAccountModel))
    if (!err) {
    this.modalRef.close()
  }
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

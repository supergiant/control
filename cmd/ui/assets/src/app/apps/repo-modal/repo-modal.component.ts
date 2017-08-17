import { Component, OnInit, AfterViewInit, OnDestroy,ViewChild, ElementRef } from '@angular/core';
import { NgbModal, ModalDismissReasons, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { Subscription } from 'rxjs/Subscription';
import { RepoModalService } from './repo-modal.service';
import { Notifications } from '../../shared/notifications/notifications.service'
import { Supergiant } from '../../shared/supergiant/supergiant.service'
import { Observable } from 'rxjs/Rx';



@Component({
  selector: 'app-repo-modal',
  templateUrl: './repo-modal.component.html',
  styleUrls: ['./repo-modal.component.css']
})
export class RepoModalComponent implements OnInit {
  private modalRef: NgbModalRef;
  private subscription: Subscription;
  private repos = [];
  @ViewChild('repoModal') content: ElementRef;


  constructor(
    private modalService: NgbModal,
    private repoModalService: RepoModalService,
    private notifications: Notifications,
    private supergiant: Supergiant,
  ) { }


  //get accouts when page loads
  ngOnInit() {
    this.getRepos()
  }
  //get accounts
  getRepos() {
    this.subscription = Observable.timer(0, 5000)
    .switchMap(() => this.supergiant.HelmRepos.get()).subscribe(
      (repos) => { this.repos = repos.items},
      () => {});
  }
  // Data init after load
  ngAfterViewInit() {
    // Check for messages from the new Cloud Accont dropdown, or edit button.
    this.subscription = this.repoModalService.newModal.subscribe( message => {
      {};
      // open the New/Edit modal
      {this.open(this.content)};});
  }

  open(content) {
    let options: NgbModalOptions = {
      size: 'lg'
    };
    this.modalRef = this.modalService.open(content, options);
  }

  ngOnDestroy(){
    this.subscription.unsubscribe();
  }

  onSubmit() {
    this.modalRef.close();
  }
}

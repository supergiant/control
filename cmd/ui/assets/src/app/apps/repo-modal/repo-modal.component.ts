import { Component, OnInit, AfterViewInit, OnDestroy, ViewChild, ElementRef } from '@angular/core';
import { NgbModal, ModalDismissReasons, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { Subscription } from 'rxjs/Subscription';
import { RepoModalService } from './repo-modal.service';
import { Notifications } from '../../shared/notifications/notifications.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Observable } from 'rxjs/Observable';

export class RepoModel {
  repo = {
    'model': {
      'name': '',
      'url': ''
    }
  };
}

@Component({
  selector: 'app-repo-modal',
  templateUrl: './repo-modal.component.html',
  styleUrls: ['./repo-modal.component.scss']
})
export class RepoModalComponent implements OnInit, AfterViewInit, OnDestroy {
  private modalRef: NgbModalRef;
  subscriptions = new Subscription();
  private repos = [];
  private new = false;
  private name: string;
  private url: string;
  private repoModel = new RepoModel;
  @ViewChild('repoModal') content: ElementRef;


  constructor(
    private modalService: NgbModal,
    private repoModalService: RepoModalService,
    private notifications: Notifications,
    private supergiant: Supergiant,
  ) { }


  ngOnInit() {
    this.getRepos();
  }

  getRepos() {
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.HelmRepos.get()).subscribe(
      (repos) => { this.repos = repos.items; },
      () => { }));
  }

  ngAfterViewInit() {
    // Check for messages from the new Cloud Accont dropdown, or edit button.
    this.subscriptions.add(this.repoModalService.newModal.subscribe(
      (message) => {
        { }
        // open the New/Edit modal
        { this.open(this.content); }
      }));
    this.subscriptions.add(this.repoModalService.repoModalResponse.subscribe(
      (message) => {
        this.repoModalService.resetSelected();
        this.new = false;
      },
      (err) => {
        this.repoModalService.resetSelected();
        this.new = false;
      },
    ));
  }

  open(content) {
    const options: NgbModalOptions = {
      size: 'lg'
    };
    this.modalRef = this.modalService.open(content, options);
    this.modalRef.result.then(
      (window) => {
        this.repoModalService.resetSelected();
        this.repoModalService.repoModalResponse.next('closed');
      },
      (err) => {
        this.repoModalService.resetSelected();
        this.repoModalService.repoModalResponse.next('closed');
      },
    );
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }
  save() {
    this.repoModel.repo.model.name = this.name;
    this.repoModel.repo.model.url = this.url;
    this.subscriptions.add(this.supergiant.HelmRepos.create(this.repoModel.repo.model).subscribe(
      (success) => {
        this.notifications.display('success', 'Repo: ' + this.name, 'Created...');
        this.new = false;
        this.getRepos();
        this.name = '';
        this.url = '';
      },
      (err) => { this.notifications.display('error', 'Create Error:', err); },
    ));
  }

  selected() {
    if (this.repoModalService.returnSelected().length > 0) {
      return true;
    } else {
      return false;
    }
  }

  deleteRepo() {
    const selectedItems = this.repoModalService.returnSelected();
    if (selectedItems.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No Repo Selected.');
    } else {
      for (const repo of selectedItems) {
        this.subscriptions.add(this.supergiant.HelmRepos.delete(repo.id).subscribe(
          (data) => {
            this.notifications.display('success', 'Repo: ' + repo.name, 'Deleted...');
            this.repoModalService.resetSelected();
            this.getRepos();
          },
          (err) => {
            this.notifications.display('error', 'Repo: ' + repo.name, 'Error:' + err);
          },
        ));
      }
    }
  }

  onSubmit() {
    this.new = true;
  }
}

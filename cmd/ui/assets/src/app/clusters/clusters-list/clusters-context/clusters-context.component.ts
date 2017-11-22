import { Component, OnInit, ViewEncapsulation, TemplateRef, ViewChild } from '@angular/core';
import { ContextMenuService, ContextMenuComponent } from 'ngx-contextmenu';

@Component({
  selector: 'app-clusters-context',
  templateUrl: './clusters-context.component.html',
  styleUrls: ['./clusters-context.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class ClustersContextComponent implements OnInit {
  @ViewChild(ContextMenuComponent) public basicMenu: ContextMenuComponent;
  constructor(private contextMenuService: ContextMenuService) { }

  ngOnInit() {
  }

}

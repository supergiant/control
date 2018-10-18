import { Component, OnInit, OnDestroy, AfterViewInit, Inject, ViewChild, ElementRef } from '@angular/core';
import { MatDialog, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';

@Component({
  selector: 'task-logs',
  templateUrl: './task-logs.component.html',
  styleUrls: ['./task-logs.component.scss']
})
export class TaskLogsComponent implements OnInit, OnDestroy, AfterViewInit {

  constructor(
    public dialogRef: MatDialogRef<TaskLogsComponent>,
      @Inject(MAT_DIALOG_DATA) public data: any
  ) { this.taskId = data.taskId }

  @ViewChild("editor") editor;

  conn: any

  taskId: any;
  logsString = "";
  logsArray = [];
  userScrolled = false

  userScrollEvent(e) { console.log(e) }

  scrollToBottom() {
    // console.log(this.editor.getEditor());
    const renderer = this.editor.getEditor().renderer
    const height = renderer.$size.scrollerHeight;
    if (renderer.scrollTop > 700) {
      this.editor.getEditor().renderer.scrollTop = this.editor.getEditor().renderer.scrollTop + ( height * 2 );
    }
  }

  updateLogs(e) {
    const newMessage = e.data + "\r\n"
    this.logsString = this.logsString + newMessage;
    if (!this.userScrolled) {
      this.scrollToBottom();
    }
  }

  openConn(taskId) {
    this.conn = new WebSocket("ws://localhost:8080/tasks/" + taskId + "/logs");
    this.conn.onmessage = e => { setTimeout(() => this.updateLogs(e), 1); }
  }

  ngOnInit() {
    this.openConn(this.taskId);
  }

  ngAfterViewInit() {
  }

  ngOnDestroy() {
    this.conn.close();
  }

}

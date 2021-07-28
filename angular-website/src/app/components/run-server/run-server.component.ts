import { Router } from '@angular/router';
import { BlockService } from './../../service/block.service';
import { Component, OnInit, NgZone } from '@angular/core';
import * as $ from 'jquery';

@Component({
  selector: 'app-run-server',
  templateUrl: './run-server.component.html',
  styleUrls: ['./run-server.component.css']
})
export class RunServerComponent implements OnInit {

  User: any;

  constructor(
    private router: Router,
    private ngZone: NgZone,
    private blockService: BlockService,
  ) { 
  }

  ngOnInit() {
    this.runServer();
  }
  
  // Function to run Satoshi Server //
  runServer() {
    this.blockService.runServer().subscribe(
    (res) => {
      if (res['status'] == 200)
      {
        console.log("Server is running...")
      }
    }, (error) => {
      console.log(error);
      console.exception("ss");
    });
  }
}

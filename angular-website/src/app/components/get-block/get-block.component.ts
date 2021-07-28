import { Router } from '@angular/router';
import { BlockService } from './../../service/block.service';
import { Component, OnInit, NgZone } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from "@angular/forms";
import * as $ from 'jquery';

@Component({
  selector: 'app-get-block',
  templateUrl: './get-block.component.html',
  styleUrls: ['./get-block.component.css']
})
export class GetBlockComponent implements OnInit {

  User: any;
  Blocks: any;
  showBlocks: any;
  submitted = false;
  playerForm: FormGroup;

  constructor(
    public fb: FormBuilder,
    private router: Router,
    private ngZone: NgZone,
    private blockService: BlockService,
  ) { 
  }

  ngOnInit() {
    this.getTheBlock();
    this.getUser();
    $(document).ready(function(){
      $('.nav_btn').click(function(){
        $('.mobile_nav_items').toggleClass('active');
      });
    });
  }
  
  getUser() {
    if (this.blockService.isLoggedIn == false)
    {
      alert("Access not allowed");
      this.router.navigateByUrl('/login')
    }
    else
    {   
      this.blockService.getUser().subscribe(
      (res) => {
        if (res['status'] == 200)
        {
          this.User = res['body']['result'];
        }
      }, (error) => {
        console.log(error);
        console.exception("ss");
      });
    }
  }

  // Function to get the specific block after linking from the CV //
  getTheBlock() {
    this.blockService.getSpecificBlock(localStorage.getItem("blockHash")).subscribe(
      (res) => {
        if (res['status'] == 200)
        {
          console.log(res);
          this.Blocks = res['body']['verifiedBlock'];
          console.log(this.Blocks);
        }
      }, (error) => {
        console.log(error);
        console.exception("ss");
    });
  }

  logout() {
    this.blockService.logout();
    this.ngZone.run(() => this.router.navigateByUrl('/login'))
  }

}

import { ActivatedRoute, Router } from '@angular/router';
import { BlockService } from './../../service/block.service';
import { Component, OnInit, NgZone } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from "@angular/forms";
import * as $ from 'jquery';

@Component({
  selector: 'app-verify-content',
  templateUrl: './verify-content.component.html',
  styleUrls: ['./verify-content.component.css']
})
export class VerifyContentComponent implements OnInit {

  User: any;

  constructor(
    public fb: FormBuilder,
    private router: Router,
    private ngZone: NgZone,
    private actRoute: ActivatedRoute,
    private blockService: BlockService,
  ) { 
  }

  ngOnInit() {
    let hash = this.actRoute.snapshot.paramMap.get('hash');
    console.log(hash);
    this.verifyHash(hash);
  }
  
  // Function to verify the block hash //
  verifyHash(hash) {
    this.blockService.mineBlock(hash).subscribe(
    (res) => {
      if (res['status'] == 200)
      {
        console.log("Content Verified!")
      }
    }, (error) => {
      console.log(error);
      console.exception("ss");
    });
  }
}

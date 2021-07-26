import { Router } from '@angular/router';
import { BlockService } from './../../service/block.service';
import { Component, OnInit, NgZone } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from "@angular/forms";

@Component({
  selector: 'app-show-blocks',
  templateUrl: './show-blocks.component.html',
  styleUrls: ['./show-blocks.component.css']
})
export class ShowBlocksComponent implements OnInit {

  Blocks: any;

  constructor(
    public fb: FormBuilder,
    private router: Router,
    private ngZone: NgZone,
    private blockService: BlockService,
  ) { 
  }

  ngOnInit() {
    this.showBlocks();
  }

  showBlocks() {
    this.blockService.showBlocks().subscribe(
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
}

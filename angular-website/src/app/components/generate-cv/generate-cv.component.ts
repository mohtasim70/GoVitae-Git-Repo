import { Router } from '@angular/router';
import { BlockService } from './../../service/block.service';
import { Component, OnInit, NgZone } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from "@angular/forms";
import { ViewChild, ElementRef } from '@angular/core';
import {jsPDF} from 'jspdf';
import html2canvas from 'html2canvas';
import * as $ from 'jquery';

@Component({
  selector: 'app-generate-cv',
  templateUrl: './generate-cv.component.html',
  styleUrls: ['./generate-cv.component.css']
})
export class GenerateCVComponent implements OnInit {

  CV: any;
  User: any;
  submitted = false;
  playerForm: FormGroup;
  @ViewChild('pdf', {static: false}) pdf: ElementRef;

  constructor(
    public fb: FormBuilder,
    private router: Router,
    private ngZone: NgZone,
    private blockService: BlockService,
  ) {
  }

  ngOnInit() {
    this.getCV();
    this.getUser();
    $(document).ready(function(){
      $('.nav_btn').click(function(){
        $('.mobile_nav_items').toggleClass('active');
      });
    });
  }

  getCV() {
    this.blockService.generateCV().subscribe(
    (res) => {
      if (res['status'] == 200)
      {
        console.log(res);
        this.CV = res['body']['cv'];
        console.log(this.CV);
      }
    }, (error) => {
      console.log(error);
      console.exception("ss");
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

  logout() {
    this.blockService.logout();
    this.ngZone.run(() => this.router.navigateByUrl('/login'))
  }

  public downloadAsPDF() {
    let data = document.getElementById('pdf');
            html2canvas(data).then(canvas => {
              
            // let pdf = new jspdf('l', 'cm', 'a4'); //Generates PDF in landscape mode
            let pdf = new jsPDF('p', 'cm', 'a4');// Generates PDF in portrait mode
            // pdf.addImage(contentDataURL, 'PNG', 0, 0, 29.7, 21.0);
            pdf.save('CV.pdf');
          });
 }

}

import { Router } from '@angular/router';
import { BlockService } from './../../service/block.service';
import { Component, OnInit, NgZone } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from "@angular/forms";
import { ViewChild, ElementRef } from '@angular/core';
import { faLink } from '@fortawesome/free-solid-svg-icons';
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
  faLink = faLink;
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
    $(document).ready(function(){
      $('[data-toggle=tooltip]').hover(function(){
          // on mouseenter
          $(this).tooltip('show');
      }, function(){
          // on mouseleave
          $(this).tooltip('hide');
      });
    });
  }

  // Function to generate CV //
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

  // Function to link to the specific block //
  getTheBlock(hash: any) {
    localStorage.setItem("blockHash", hash)
    this.ngZone.run(() => this.router.navigateByUrl('/getBlock'))
  }
 
  logout() {
    this.blockService.logout();
    this.ngZone.run(() => this.router.navigateByUrl('/login'))
  }

  // Function to generate PDF of CV //
  downloadAsPDF() {
    
    document.getElementById('carder').style.boxShadow = 'none';

    document.getElementById('buttonAdd').style.display = 'none';

    let DATA = document.getElementById('pdf');
      
    html2canvas(DATA).then(canvas => {
        
        let fileWidth = 208;
        let fileHeight = canvas.height * fileWidth / canvas.width;
        
        const FILEURI = canvas.toDataURL('image/png')
        let PDF = new jsPDF('p', 'mm', 'a4');
        let position = 0;
        PDF.addImage(FILEURI, 'PNG', 0, position, fileWidth, fileHeight)
        
        PDF.save('CV.pdf');
    });     

    document.getElementById('carder').style.boxShadow = '0px 0px 30px rgba(0, 0, 0, 0.1)';

    document.getElementById('buttonAdd').style.display = 'block';
  }

}

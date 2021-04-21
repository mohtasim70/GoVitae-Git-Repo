import { Injectable } from '@angular/core';
import { Observable, throwError } from 'rxjs';
import { catchError, map, timeout } from 'rxjs/operators';
import { HttpClient, HttpHeaders, HttpErrorResponse,  HttpParams, HttpRequest } from '@angular/common/http';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root'
})
export class BlockService {

  mainURL:string = "";

  constructor(private http: HttpClient, public router: Router) { }

  // Register Function //

  register(data): Observable<any> {
    let url = `${this.mainURL}/registerUser`;

    const req = new HttpRequest('POST', url, data);
    return this.http.request(req);
  }

  // Login Function //

  login(data): Observable<any> {
    let url = `${this.mainURL}/loginUser`;

    const req = new HttpRequest('POST', url, data);
    return this.http.request(req);
  }

  // Logout Function //

  logout() {
    let removeToken = localStorage.removeItem("TOKEN");
    if (removeToken == null) {
      this.router.navigate(['login']);
    }
  }

  // Add Course Function //

  addCourse(data): Observable<any> {
    let url = `${this.mainURL}/addCourseUser`;

    const req = new HttpRequest('POST', url, data);
    return this.http.request(req);
  }

  // Add Project Function //

  addProject(data): Observable<any> {
    let url = `${this.mainURL}/addProjectUser`;

    const req = new HttpRequest('POST', url, data);
    return this.http.request(req);
  }

  // Generate CV Function //

  generateCV(): Observable<any> {
    let url = `${this.mainURL}/generateCVUser`;

    const req = new HttpRequest('GET', url);
    return this.http.request(req);
  }

  // Get Unverified Blocks Function //

  getBlocks(): Observable<any> {
    let url = `${this.mainURL}/getBlocksUser`;

    const req = new HttpRequest('GET', url);
    return this.http.request(req);
  }
  
  // Get User Function //

  getUser() {
    let url = `${this.mainURL}/getUser`;

    const req = new HttpRequest('GET', url);
    return this.http.request(req);
  }

  // Get All Users Function //

  getAllUsers() {
    let url = `${this.mainURL}/getAllUsers`;

    const req = new HttpRequest('GET', url);
    return this.http.request(req);
  }

  // Get Verified Content Function //

  getVerifiedContent() {
    let url = `${this.mainURL}/getVerifyContent`;

    const req = new HttpRequest('GET', url);
    return this.http.request(req);
  }

  // Get Verified CVs Function //

  getVerifiedCVs(data) {
    let url = `${this.mainURL}/getVerifiedCVs`;

    const req = new HttpRequest('POST', url, data);
    return this.http.request(req);
  }

  // Mine Block Function //

  mineBlock(hash) {
    let url = `${this.mainURL}/mineBlockMiner/${hash}`;

    const req = new HttpRequest('GET', url);
    return this.http.request(req);
  }

  // Get Logged In User //

  get isLoggedIn(): boolean {
    let authToken = localStorage.getItem("TOKEN");
    return (authToken !== null) ? true : false;
  }
  
}

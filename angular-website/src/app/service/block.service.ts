import { Injectable } from '@angular/core';
import { Observable, throwError } from 'rxjs';
import { catchError, map, timeout } from 'rxjs/operators';
import { HttpClient, HttpHeaders, HttpErrorResponse,  HttpParams, HttpRequest } from '@angular/common/http';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root'
})
export class BlockService {

  adminURL:string = "http://localhost:3334";

  constructor(private http: HttpClient, public router: Router) { }

  // Register Function //

  register(data): Observable<any> {
    let url = `${this.adminURL}/register`;

    const req = new HttpRequest('POST', url, data);
    return this.http.request(req);
  }

  // Login Function //

  login(data): Observable<any> {
    let url = `${this.adminURL}/login`;

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
    let url = `${this.adminURL}/addCourse`;

    const req = new HttpRequest('POST', url, data);
    return this.http.request(req);
  }

  // Add Project Function //

  addProject(data): Observable<any> {
    let url = `${this.adminURL}/addProject`;

    const req = new HttpRequest('POST', url, data);
    return this.http.request(req);
  }

  // Generate CV Function //

  generateCV(): Observable<any> {
    let url = `${this.adminURL}/generateCV`;

    const req = new HttpRequest('GET', url);
    return this.http.request(req);
  }

  // Get Unverified Blocks Function //

  getBlocks(): Observable<any> {
    let url = `${this.adminURL}/getBlocks`;

    const req = new HttpRequest('GET', url);
    return this.http.request(req);
  }
  
  // Get User Function //

  getUser() {
    let url = `${this.adminURL}/getUser`;

    const req = new HttpRequest('GET', url);
    return this.http.request(req);
  }

  // Get All Users Function //

  getAllUsers() {
    let url = `${this.adminURL}/getAllUsers`;

    const req = new HttpRequest('GET', url);
    return this.http.request(req);
  }

  // Get Logged In User //

  get isLoggedIn(): boolean {
    let authToken = localStorage.getItem("TOKEN");
    return (authToken !== null) ? true : false;
  }
  
}

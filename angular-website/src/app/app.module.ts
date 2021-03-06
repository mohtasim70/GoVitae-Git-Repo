import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';

import {AutocompleteLibModule} from 'angular-ng-autocomplete';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { AddCourseComponent } from './components/add-course/add-course.component';
import { LoginComponent } from './components/login/login.component';
import { RegisterComponent } from './components/register/register.component';
import { AddProjectComponent } from './components/add-project/add-project.component';
import { DashboardComponent } from './components/dashboard/dashboard.component';
import { UnverifiedBlocksComponent } from './components/unverified-blocks/unverified-blocks.component';
import { GenerateCVComponent } from './components/generate-cv/generate-cv.component';
import { HomeComponent } from './components/home/home.component';
import { SearchCVComponent } from './components/search-cv/search-cv.component';
import { VerifyContentComponent } from './components/verify-content/verify-content.component';
import { RunServerComponent } from './components/run-server/run-server.component';
import { GetBlockComponent } from './components/get-block/get-block.component';
import { ShowBlocksComponent } from './components/show-blocks/show-blocks.component';

@NgModule({
  declarations: [
    AppComponent,
    AddCourseComponent,
    LoginComponent,
    RegisterComponent,
    AddProjectComponent,
    DashboardComponent,
    UnverifiedBlocksComponent,
    GenerateCVComponent,
    HomeComponent,
    SearchCVComponent,
    VerifyContentComponent,
    RunServerComponent,
    GetBlockComponent,
    ShowBlocksComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FormsModule,
    ReactiveFormsModule,
    HttpClientModule,
    AutocompleteLibModule,
    FontAwesomeModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }

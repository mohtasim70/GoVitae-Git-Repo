import React, { Component } from 'react';
import { Link } from "react-router-dom";
import { Router, Route, browserHistory, IndexRoute } from 'react-router'
// import "../Components/homepage.css"

//import the testimonials
// import Testimonials from "../components/testimonials";


class Home extends Component {
  render() {
    return (
      <div className="main">
          <header>
         
              <div className="header-data">
                <h1>GoVitae</h1>
                <p>Integer tincidunt tempor libero in condimentum.</p>
              </div>  
          </header>

          <section className="container">
            <div className="cards">
                <div>
                   <div className="icons">
                   <i className="fa fa-cog" aria-hidden="true"></i>
                   </div>
                    <div className="content">
                        <p className="content-title">Etiam vitae</p>
                        <p>Etiam vitae urna gravida, fringilla dolor tristique, pellentesque arcu. Proin finibus dignissim dui, hendrerit fringilla ante elementum sed. Nullam facilisis urna quam, a efficitur felis ullamcorper in.</p>

                    </div>
                </div>
            </div>

            <div className="cards">
                <div>
                   <div className="icons">
                   <i className="fa fa-calendar" aria-hidden="true"></i>
                   </div>
                    <div className="content">
                        <p className="content-title">Aenean quis</p>
                        <p>Etiam vitae urna gravida, fringilla dolor tristique, pellentesque arcu. Proin finibus dignissim dui, hendrerit fringilla ante elementum sed. Nullam facilisis urna quam, a efficitur felis ullamcorper in.</p>

                    </div>
                </div>
            </div>
            <div className="cards">
                <div>
                   <div className="icons">
                   <i className="fa fa-building" aria-hidden="true"></i>
                   </div>
                    <div className="content">
                        <p className="content-title">Praesent varius</p>
                        <p>Etiam vitae urna gravida, fringilla dolor tristique, pellentesque arcu. Proin finibus dignissim dui, hendrerit fringilla ante elementum sed. Nullam facilisis urna quam, a efficitur felis ullamcorper in.</p>

                    </div>
                </div>
            </div>

          </section>


          <section className="dual-card">
              <div>
                  <p className="title">Cras sollicitudin turpis</p>
                  <p>Phasellus cursus dictum tortor vel imperdiet. Nunc at lacus sapien. Etiam ipsum quam, tincidunt ut lobortis et, finibus vitae dui. Vestibulum accumsan magna a metus vehicula, sit amet lobortis lectus finibus. Integer dapibus, libero sed suscipit egestas, sem magna porta nibh, in lobortis mi ipsum scelerisque turpis.</p>
              </div>
              <div>
                  <img src="https://images.unsplash.com/photo-1506703248799-a1b1896290e8?auto=format&fit=crop&w=634&q=80" />
              </div>
          </section>

          <section className="img-section">
              <p>Image section</p>
          </section>
      
        <section className="imgs-bar">
            <img src="https://images.unsplash.com/photo-1496284045406-d3e0b918d7ba?auto=format&fit=crop&w=2134&q=80" />
            <img src="https://images.unsplash.com/photo-1500462918059-b1a0cb512f1d?auto=format&fit=crop&w=634&q=80" />
            <img src="https://images.unsplash.com/photo-1504420379316-1a01ae921844?auto=format&fit=crop&w=634&q=80" />
            <img src="https://images.unsplash.com/photo-1496979149951-4fd28cf221e7?auto=format&fit=crop&w=634&q=80" />
        </section>

        {/* <Testimonials /> */}


        <footer>
            <div className="data">
                <ul>
                    <li><p className="title">LOGO</p></li>
                </ul>
                <ul>
                    <li><a href="#">Home</a></li>
                    <li><a href="#">About Us</a></li>
                    <li><a href="#">Services</a></li>
                </ul>
            </div>
            <div className="media">
                <ul>
                    <li><a href="#"><i className="fa fa-facebook" aria-hidden="true"></i></a></li>
                    <li><a href="#"><i className="fa fa-twitter" aria-hidden="true"></i></a></li>
                    <li><a href="#"><i className="fa fa-twitch" aria-hidden="true"></i></a></li>
                </ul>
            </div>
        </footer>
     
      </div>

      
    );
  }
}

export default Home;
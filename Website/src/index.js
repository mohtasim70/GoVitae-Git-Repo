import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import Blockchain from './Blockchain'

import reportWebVitals from './reportWebVitals';
import { Route, Link, BrowserRouter as Router } from 'react-router-dom'

const routing = (
  <Router>
       <header>
      <nav>
      <ul>
        <li>
          <Link to="/">Home</Link>
        </li>
        <li>
          <Link to="/create">Add Verified Course</Link>
        </li>
      </ul>
      <Route exact path="/" component={App} />
      <Route path="/create" component={Blockchain}/>
      </nav>
      </header>
  </Router>
)

ReactDOM.render(
  // <React.StrictMode>
  //   <App />
    
  // </React.StrictMode>,
  routing,
  document.getElementById('root')
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();

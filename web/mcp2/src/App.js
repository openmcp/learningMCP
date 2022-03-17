// import logo from './logo.svg';
import React from 'react';
import { Route } from 'react-router-dom';
import Home from './page/Home';
import './App.css';

const App = () => {
  return (
    <>
      <Route path="/" component={Home} exact={true} />
    </>
  );
}

export default App;

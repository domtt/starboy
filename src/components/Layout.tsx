import React from 'react';
import "./Layout.css"
import {useToken} from '../state/token';

const Layout = ({children}: {children: any}) => {
  const token = useToken();
  const signIn = () => {
    window.location.href = "http://localhost:3000/auth/github";
  }
  return <div>
    <div className="toolbar">
      <div className="title">Starboy</div>

      {!token && <button onClick={signIn} className="flat">Sign In with GitHub</button>}
    </div>
    <div className="container">
    {children}
    </div>
  </div>
}

export default Layout;

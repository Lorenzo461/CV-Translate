import React from 'react'
import 'bootstrap/dist/css/bootstrap.css';

export default function Card() {
  return (
    <div className="card" style={{width:'18rem'}}>
        
        <div className="card-body">
            <h5 className="card-title">Card title</h5>
            <p className="card-text">Some quick example text to build on the card title and make up the bulk of the card's content.</p>
            <button className="btn btn-primary">Go somewhere</button>
        </div>
    </div>
  )
}

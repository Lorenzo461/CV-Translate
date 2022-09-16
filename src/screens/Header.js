import React from 'react'
import '../css/Header.css';
import strada from '../img/strada.png';

export default function Header() {
  return (
    <div className='container-header'>
         <h1> CV Translate</h1>
         <h2>Tutte le strade portano al tuo nuovo lavoro!</h2>
        <div >
        <img src={strada} alt="tranlsate" />
        </div>
    </div>
  )
}

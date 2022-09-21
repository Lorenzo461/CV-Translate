import React from 'react'
import '../css/Header.css';
import strada from '../img/strada.png';
import Button from '../components/Button';

export default function Header() {
  return (
    <div className='container-header'>
         <h1> CV Translate</h1>
         <h2>Tutte le strade portano al tuo nuovo lavoro!</h2>
        <div >
        <img src={strada} alt="tranlsate" />
        </div>
        <div>
        <p>
          Con  <h3 className='deco'>CV</h3> <h3 className='deco'>Translate</h3> puoi tradurre velocemente e in pochi semplici passi il tuo <span> Curriculum vitae. </span><br />
          Il servizio che mettiamo a disposizione è pensato per aiutare tutte quelle persone che hanno bisogno di ottenere il proprio cv in un'altra lingua.<br />
          Cosa aspetti ? Aumenta le tue possibilità di successo con un solo click !
        </p>
        </div>
        <div>
        <Button text="Traduci il tuo CV" />
        </div>
    </div>
  )
}

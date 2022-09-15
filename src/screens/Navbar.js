import React from 'react'
import Menu from '../components/Menu'
import '../css/Navbar.css'

export default function Navbar() {
  return (
    <div className='navbar'>
        <div>
            <h1>
                CV Translate
            </h1>
        </div>
        <div>
        <Menu />
        </div>
        
    </div>
  )
}

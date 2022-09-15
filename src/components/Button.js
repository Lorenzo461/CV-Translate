import React from 'react'
import './Button.css';

export default function Button(props) {
  return (
    <button className='btn'  onClick={props.handleClick} type="button"><span>{props.text}</span></button>
  )
}

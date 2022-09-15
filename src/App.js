import Navbar from './screens/Navbar';
import Button from './components/Button';
import logo from './img/log.png';
import './App.css'

function App() {
  return (
    <div>
      <Navbar />
      <div className='container'>
        <div className='uno'>
          <h2>
          Traduci online il tuo CV professionale <br />
          dai una svolta alla tua vita.<br />
          Lavora in tutto il mondo senza limiti!
          </h2>
          <Button text="Traduci CV"/>
        </div>
        <div className='due'>
          <img src={logo} alt="tranlsate" />
        </div>
      </div>
    </div>
  );
}

export default App;

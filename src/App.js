import Navbar from './screens/Navbar';
import Header from './screens/Header';
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
          <Button text="Traduci il tuo CV"/>
        </div>
        <div className='due'>
          <img src={logo} alt="tranlsate" />
        </div>
      </div>
      <Header/>
    </div>
  );
}

export default App;

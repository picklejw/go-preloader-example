import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import './App.css';
import Item from './Item';

function Home() {
  console.log("Home Render...")
  return (
    <div className="App">
      <header className="App-header">
        <h1>This is /</h1>
        <nav>
          <Link to="/item?id=123" className="App-link">View Item 123</Link>
        </nav>
      </header>
    </div>
  );
}

function App() {
  console.log("App render...")
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/item" element={<Item />} />
      </Routes>
    </Router>
  );
}

export default App;

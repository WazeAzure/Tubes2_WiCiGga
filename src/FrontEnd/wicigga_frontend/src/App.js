import './App.css';
import { useState, useEffect } from 'react';


function App() {
  const [data, setData] = useState([]);
  // testing using state

  useEffect(() => {
    fetch('http://localhost:4000/path')
      .then((res) => {
        return res.json();
      }).then((data) => {
        console.log(data);
        setData(data.path);
      })
  }, [])

  //Search Bar Handler
  const [value1, setValue1] = useState("");
  const [value2, setValue2] = useState("");

  //Toggle Button Handler, false = BFS, true = IDS
  const [buttonState, setButtonState] = useState(false);


  const onChange1 = (event) => {
    setValue1(event.target.value);
  }

  const onChange2 = (event) => {
    setValue2(event.target.value);
  }

  const onSearch = (searchTerm) => {
    //fetch api
    console.log("search", searchTerm);
  }

  return (
    <div className="Head">
      <h1 className='title'>Wicigga</h1>

      {/* Left Search Section */}

      <div className='search-section'>
        <div className='search-left'>
          <div className='search-bar-container'>
            <input type="text" className='search-bar' value={value1} onChange={onChange1} />
          </div>
          <div className='dropdown-offset'>
            <div className={'dropdown' + (data.some(item => {
              const searchTerm = value1.toLowerCase();
              const pathString = item.toLowerCase();
              return searchTerm && pathString.startsWith(searchTerm) && searchTerm !== pathString;
            }) ? '' : 'dummy')}>
              {data.filter(item => {
                const searchTerm = value1.toLowerCase();
                const pathString = item.toLowerCase();

                return searchTerm && pathString.startsWith(searchTerm) && searchTerm !== pathString;
              })
                .map((item) => (
                  <li className="search-result">{item}</li>
                ))
              }
            </div>
          </div>
        </div>

        {/* Right Search Section */}

        <div className='search-right'>
          <div className='search-bar-container'>
            <input type="text" className='search-bar' value={value2} onChange={onChange2} />
          </div>
          <div className='dropdown-offset'>
            <div className={'dropdown' + (data.some(item => {
              const searchTerm = value2.toLowerCase();
              const pathString = item.toLowerCase();
              return searchTerm && pathString.startsWith(searchTerm) && searchTerm !== pathString;
            }) ? '' : 'dummy')}>
              {data.filter(item => {
                const searchTerm = value2.toLowerCase();
                const pathString = item.toLowerCase();

                return searchTerm && pathString.startsWith(searchTerm) && searchTerm !== pathString;
              })
                .map((item) => (
                  <li className="search-result">{item}</li>
                ))
              }
            </div>
          </div>
        </div>
      </div>
      <div className='box-1'>
        <div className='button-mode' onClick={() => setButtonState(!buttonState)}>
          {buttonState ? <p>IDS</p> : <p>BFS</p>}
        </div>
        <div className='button-search' onClick={() => onSearch()}>Search</div>
      </div>
    </div>
  );
}

export default App;

{/* <div className='search-container'>
        <div className='search-bar-outer'>
          <div className='search-bar'>
            <input type="text" value={value1} onChange={onChange1} />
          </div>
          <div className={'dropdown' + (data.some(item => {
            const searchTerm = value1.toLowerCase();
            const pathString = item.toLowerCase();
            return searchTerm && pathString.startsWith(searchTerm) && searchTerm !== pathString;
          }) ? '' : '-empty')}>
            {data.filter(item => {
              const searchTerm = value1.toLowerCase();
              const pathString = item.toLowerCase();

              return searchTerm && pathString.startsWith(searchTerm) && searchTerm !== pathString;
            })
              .map((item) => (
                <li className="search-result">{item}</li>
              ))
            }
          </div>
        </div>
      </div> */}
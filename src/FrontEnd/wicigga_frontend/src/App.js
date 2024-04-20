import './App.css';
import { useState, useEffect } from 'react';


function App() {
  //Autocomplete data
  const [data1, setData1] = useState([]);
  const [data2, setData2] = useState([]);

  //Search Bar Handler
  const [value1, setValue1] = useState("");
  const [value2, setValue2] = useState("");

  //Toggle Button Handler, false = BFS, true = IDS
  const [buttonState, setButtonState] = useState(false);

  //Input handler
  const onChange1 = (event) => {
    setValue1(event.target.value);
  }

  const onChange2 = (event) => {
    setValue2(event.target.value);
  }

  //Fetch data for autocomplet from wikipedia's API
  useEffect(() => {
    if (value1 !== "") {
      fetch("https://en.wikipedia.org/w/api.php?action=query&origin=*&prop=extracts&format=json&formatversion=2&list=search&srsearch=" + value1)
        .then((res) => {
          return res.json();
        }).then((jsonDat) => {
          console.log(jsonDat.query);
          return jsonDat.query;
        }).then((que) => {
          console.log(que.search.map(item => item.title));
          setData1(que.search.map(item => item.title));;
        });
    }
  }, [value1])

  useEffect(() => {
    if (value2 !== "") {
      fetch("https://en.wikipedia.org/w/api.php?action=query&origin=*&prop=extracts&format=json&formatversion=2&list=search&srsearch=" + value2)
        .then((res) => {
          return res.json();
        }).then((jsonDat) => {
          console.log(jsonDat.query);
          return jsonDat.query;
        }).then((que) => {
          console.log(que.search.map(item => item.title));
          setData2(que.search.map(item => item.title));;
        });
    }
  }, [value2])

  //Search Button Handler
  const onSearch = (searchTerm) => {
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
            <div className={'dropdown' + (data1.some(item => {
              const searchTerm = value1.toLowerCase();
              const pathString = item.toLowerCase();
              return searchTerm && pathString.startsWith(searchTerm);
            }) ? '' : 'dummy')}>
              {data1.filter(item => {
                const searchTerm = value1.toLowerCase();
                const pathString = item.toLowerCase();

                return searchTerm && pathString.startsWith(searchTerm);
              }).slice(0, 5)
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
            <div className={'dropdown' + (data2.some(item => {
              const searchTerm = value2.toLowerCase();
              const pathString = item.toLowerCase();
              return searchTerm && pathString.startsWith(searchTerm) && searchTerm !== pathString;
            }) ? '' : 'dummy')}>
              {data2.filter(item => {
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
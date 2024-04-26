import './App.css';
import { useState, useEffect } from 'react';
import GraphShow from './components/Graph/GraphComponent';


function App() {
  //Autocomplete data
  const [data1, setData1] = useState([]);
  const [data2, setData2] = useState([]);

  //Search bar value
  const [value1, setValue1] = useState("");
  const [value2, setValue2] = useState("");

  //Toggle button handler, false = BFS, true = IDS
  const [buttonState, setButtonState] = useState(false);
  const [selected1, setSelected1] = useState(false);
  const [selected2, setSelected2] = useState(false);

  //Data completion handler
  const [dataComplete, setDataComplete] = useState(false)
  const [form1, setFrom1Red] = useState(false)
  const [form2, setFrom2Red] = useState(false)

  // Data for GRAPH
  const [nodes, setNodes] = useState([])
  const [edges, setEdges] = useState([])

  // hasil dari API
  // durasi eksekusi program
  const [duration, setDuration] = useState(0)
  // boolean. jika path ditemukan dari start -> end
  const [pathFound, setPathFound] = useState(false)
  // pesan jika path tidak ditemukan / error lainnya
  const [message, setMessage] = useState("")

  //Search input handler
  const onChange1 = (event) => {
    setValue1(event.target.value);
    setSelected1(false);
    setFrom1Red(false);
  }

  const onChange2 = (event) => {
    setValue2(event.target.value);
    setSelected2(false);
    setFrom2Red(false);
  }

  const turnFormRed = (data1, data2) => {
    if (data1.length === 0 || value1 === '') {
      setFrom1Red(true);
      setValue1('');
    }
    if (data2.length === 0 || value2 === '') {
      setFrom2Red(true);
      setValue2('');
    }
  }

  useEffect(() => {
    setDataComplete(data1.length !== 0 && data2.length !== 0 && value1 !== '' && value2 !== '');
  }, [data1, data2, value1, value2]);



  //Fetch data for autocomplete from wikipedia's API

  useEffect(() => {
    if (value1 !== "") {
      fetch("https://en.wikipedia.org/w/api.php?action=query&format=json&formatversion=2&origin=*&list=search&srsearch=" + encodeURIComponent(value1))
        .then((res) => {
          return res.json();
        }).then((jsonDat) => {
          // console.log(jsonDat.query);
          return jsonDat.query;
        }).then((que) => {
          // console.log(que.search.map(item => item.title));
          setData1(que.search.map(item => item.title));;
        });
    }
  }, [value1])

  useEffect(() => {
    if (value2 !== "") {
      fetch("https://en.wikipedia.org/w/api.php?action=query&format=json&formatversion=2&origin=*&list=search&srsearch=" + encodeURIComponent(value2))
        .then((res) => {
          return res.json();
        }).then((jsonDat) => {
          // console.log(jsonDat.query);
          return jsonDat.query;
        }).then((que) => {
          // console.log(que.search.map(item => item.title));
          setData2(que.search.map(item => item.title));;
        });
    }
  }, [value2])

  // fetch api to backend , post data to API

  const sendData = async (dataToSend) => {
    // console.log(value1 === '')
    // console.log(value2 === '')
    try {
      fetch('http://localhost:4000/api', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(dataToSend)
      })
        .then((response) => {
          if (!response.ok) {
            console.error('Failed to send data:', response.statusText);
          }

          return response.json();
        })
        .then((data) => {
          setDuration(data.time);
          setMessage(data.message);
          setPathFound(data.status);
          setEdges(data.edges);
          setNodes(data.nodes);

          // console.log(data)

        })
    } catch (err) {
      console.log(err)
    }

  }

  //Search button handler

  const onSearch = () => {

    const dataToSend = {
      start: value1,
      end: value2,
      method: buttonState ? 'IDS' : 'BFS'
    }
    if (dataComplete) {

      console.log(dataToSend)
      sendData(dataToSend)
    }
  }



  return (
    <div className="Head">
      <h1 className='title'>Wicigga</h1>


      <div className='search-section'>

        {/* Left Search Section */}

        <div className='search-left'>
          <div className='search-bar-container'>
            {!form1 ? <input type="text" placeholder='Type here to search..' className='search-bar' value={value1} onChange={onChange1} />
              : <input type="text" placeholder='Please input form accordingly!' className='search-bar' value={value1} onChange={onChange1} />}
          </div>
          <div className='dropdown-offset'>
            <div className={'dropdown' + (data1.some(item => {
              const searchTerm = value1.toLowerCase();
              return !selected1 && searchTerm;
            }) ? '' : 'dummy')}>
              {data1.filter(item => {
                const searchTerm = value1.toLowerCase();

                // console.log(selected2)

                return !selected1 && searchTerm;
              }).slice(0, 5)
                .map((item) => (
                  <li className="search-result" onClick={() => { setValue1(item); setSelected1(true); }}>{item}</li>
                ))
              }
            </div>
          </div>
        </div>

        {/* Right Search Section */}

        <div className='search-right'>
          <div className='search-bar-container'>
            {!form2 ? <input type="text" placeholder='Type here to search..' className='search-bar' value={value2} onChange={onChange2} />
              : <input type="text" placeholder='Please input form accordingly!' className='search-bar' value={value2} onChange={onChange2} />}
          </div>
          <div className='dropdown-offset'>
            <div className={'dropdown' + (data2.some(item => {
              const searchTerm = value2.toLowerCase();
              return !selected2 && searchTerm;
            }) ? '' : 'dummy')}>
              {data2.filter(item => {
                const searchTerm = value2.toLowerCase();

                // console.log(selected2)

                return !selected2 && searchTerm;
              }).slice(0, 5)
                .map((item) => (
                  <li className="search-result" onClick={() => { setValue2(item); setSelected2(true); }}>{item}</li>
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
        {dataComplete ? <div className='button-search' onClick={() => onSearch()}>Search</div>
          : <div className='button-search' onClick={() => turnFormRed(data1, data2)}>Search</div>}
      </div>

      {/* SHOW RESULT */}
      <div>
        {
          pathFound &&
          <>
            <p>{message}</p>
          </>
        }
        {
          !pathFound &&
          <p>Path NOT Founded</p>
        }
        <p>execution time {duration}</p>
      </div>
      {/* VISUALIZE RESULT IN GRAPH */}
      <GraphShow node_list={nodes} edge_list={edges} />
    </div>
  );
}

export default App;
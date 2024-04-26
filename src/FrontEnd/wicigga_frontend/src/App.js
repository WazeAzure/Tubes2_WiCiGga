import './App.css';
import { useState, useEffect } from 'react';
import GraphShow from './components/Graph/GraphComponent';


function App() {
  //Autocomplete data
  const [data1, setData1] = useState([]);
  const [data2, setData2] = useState([]);

  //Image URls
  const [imageURLs1, setImageURLs1] = useState([]);
  const [imageURLs2, setImageURLs2] = useState([]);

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

  // const getImage = async (item, id) => {
  //   let url = "";
  //   if (item.length !== "") {
  //     try {
  //       const response = await fetch(`https://en.wikipedia.org/w/api.php?action=query&prop=pageimages&format=json&origin=*&pithumbsize=100&pageids=${id}`);
  //       const jsonDat = await response.json();
  //       const pages = jsonDat.query?.pages; // Optional chaining to prevent errors
  //       const cont = pages ? pages[id] : null;
  //       if (cont && cont.thumbnail) {
  //         url = cont.thumbnail.source;
  //         console.log(url);
  //       }
  //     } catch (error) {
  //       console.error("Error fetching image:", error);
  //     }
  //   }
  //   return url;
  // };

  useEffect(() => {
    if (data1.length !== 0) {
      Promise.all(data1.map((item) =>
        fetch(`https://en.wikipedia.org/w/api.php?action=query&prop=pageimages&format=json&origin=*&pithumbsize=100&pageids=${item.pageid}`)
          .then((res) => res.json())
          .then((jsonDat) => jsonDat.query)
          .then((queryResult) => {
            // Extract the image URL from the query result and return it
            const ImageURL = queryResult.pages[item.pageid]?.thumbnail?.source;

            return ImageURL ? ImageURL : '';
          })
      ))
        .then((imageURLs) => {
          // Set the image URLs using setImageURLs1
          setImageURLs1(imageURLs);
        })
        .catch((error) => {
          console.error("Error fetching data:", error);
        });
    }
  }, [data1]);

  useEffect(() => {
    if (data2.length !== 0) {
      Promise.all(data2.map((item) =>
        fetch(`https://en.wikipedia.org/w/api.php?action=query&prop=pageimages&format=json&origin=*&pithumbsize=100&pageids=${item.pageid}`)
          .then((res) => res.json())
          .then((jsonDat) => jsonDat.query)
          .then((queryResult) => {
            // Extract the image URL from the query result and return it
            const ImageURL = queryResult.pages[item.pageid]?.thumbnail?.source;

            return ImageURL ? ImageURL : '';
          })
      ))
        .then((imageURLs) => {
          // Set the image URLs using setImageURLs1
          setImageURLs2(imageURLs);
        })
        .catch((error) => {
          console.error("Error fetching data:", error);
        });
    }
  }, [data2]);

  // Fetch data for image search thumbnail
  // useEffect((id) => {
  //   if (data1.length !== 0) {
  //     fetch(`https://en.wikipedia.org/w/api.php?action=query&prop=pageimages&format=json&origin=*&pithumbsize=100&pageids=` + id)
  //       .then((res) => {
  //         return res.json();
  //       }).then
  //   }
  // }, [])


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
          setData1(que.search.map(item => ({ title: item.title, pageid: item.pageid })));
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
          setData2(que.search.map(item => ({ title: item.title, pageid: item.pageid })));
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

      // console.log(dataToSend)
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
            {!form1 ? <input type="text" placeholder='Type here to search..' className='search-bar' value={value1} onChange={onChange1} onBlur={() => setSelected1(true)} onFocus={() => setSelected1(false)} />
              : <input type="text" placeholder='Please input form accordingly!' className='search-bar' value={value1} onChange={onChange1} />}
          </div>
          <div className='dropdown-offset'>
            <div className={'dropdown' + (data1.some(item => {
              const searchTerm = value1.toLowerCase();
              return !selected1 && searchTerm;
            }) ? '' : '-dummy')}>
              {data1.slice(0, 5)
                .map((item, index) => (
                  <li key={index} className="search-result" onClick={() => { setValue1(item); setSelected1(true); }}>
                    <div className='img-box'>
                      <div className='image-result'>
                        <img src={imageURLs1[index]} className='image'></img>
                      </div>
                    </div>
                    <div className='result-box'>{item.title}</div>
                  </li>
                ))
              }
            </div>
          </div>
        </div>

        {/* Right Search Section */}

        <div className='search-right'>
          <div className='search-bar-container'>
            {!form2 ? <input type="text" placeholder='Type here to search..' className='search-bar' value={value2} onChange={onChange2} onBlur={() => setSelected2(true)} onFocus={() => setSelected2(false)} />
              : <input type="text" placeholder='Please input form accordingly!' className='search-bar' value={value2} onChange={onChange2} />}
          </div>
          <div className='dropdown-offset'>
            <div className={'dropdown' + (data2.some(item => {
              const searchTerm = value2.toLowerCase();
              return !selected2 && searchTerm;
            }) ? '' : '-dummy')}>
              {data2.slice(0, 5)
                .map((item, index) => (
                  <li key={index} className="search-result" onClick={() => { setValue2(item); setSelected2(true); }}>
                    <div className='img-box'>
                      <div className='image-result'>
                        <img src={imageURLs2[index]} className='image'></img>
                      </div>
                    </div>
                    <div className='result-box'>{item.title}</div>
                  </li>
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
        {/* {pathFound? <p>{message}</p>:<p>Path NOT Founded</p>} */}
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
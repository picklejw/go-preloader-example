import { useSearchParams } from 'react-router-dom';
import PreloaderAPILoader from './utils/PreloaderAPILoader';
import { useEffect, useState } from 'react';


PreloaderAPILoader('/api')

function Item() {
  const [searchParams] = useSearchParams();
  const id = searchParams.get('id');
  const [item, setItem] = useState({})
  const [err, setErr] = useState({})

  useEffect(() => {
    (async () => {
      try {
        console.log("Calling item for api data...")
        const data = await PreloaderAPILoader().get("/item?id=123"); // make sure this works, when accessing cache
        // should be for path and query parameter as a valid preloaded state
        setItem(data)
      } catch (error) {
        setErr(`Could not set item desc: ${error}`)
        console.error(error)
      }
    })()
  }, [])

          console.log("itme render..")

  return (
    <div className="App">
      <header className="App-header">
        <h1>This is item with id: {id}</h1>
        <h3>The item name is : {item.name}</h3>
      </header>
    </div>
  );
}

export default Item;

import * as react from "react"
import { indexContents } from "../contents/index"

function IndexPopup() {
  return (
    // <div
    //   style={{
    //     padding: 16
    //   }}>
    //   <h2>
    //     Welcome to your{" "}
    //     <a href="https://www.plasmo.com" target="_blank">
    //       Plasmo
    //     </a>{" "}
    //     Extension!
    //   </h2>
    //   <input onChange={(e) => setData(e.target.value)} value={data} />
    //   <a href="https://docs.plasmo.com" target="_blank">
    //     View Docs
    //   </a>
    // </div>
    <div style={{ width: '200px', height: '100px' }}>
      <button onClick={indexContents}>回答生成</button>
    </div>
  )
}

export default IndexPopup

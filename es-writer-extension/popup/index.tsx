import React from "react"
import indexContents from "./indexContents"

function IndexPopup() {
  return (
    <div style={{ width: '150px', height: '75px' }}>
      <button onClick={indexContents}>回答生成</button>
    </div>
  )
}

export default IndexPopup

import React from "react"
import { indexContents } from "../contents/index"

function IndexPopup() {
  const handleButtonClick = () => {
    console.log("Button clicked") // ボタンがクリックされたことを確認するためのログ
    indexContents()
  }

  return (
    <div style={{ width: '150px', height: '75px' }}>
      <button onClick={handleButtonClick}>回答生成</button>
    </div>
  )
}

export default IndexPopup

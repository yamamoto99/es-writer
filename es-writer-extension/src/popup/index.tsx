import { MemoryRouter } from "react-router-dom"

import Routing from "./routes/App"
import React from "react"

function IndexPopup() {
  return (
    <MemoryRouter>
      <Routing />
    </MemoryRouter>
  )
}

export default IndexPopup
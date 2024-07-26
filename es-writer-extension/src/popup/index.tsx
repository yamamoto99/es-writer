import React from "react"
import { MemoryRouter } from "react-router-dom"

import Routing from "./routes/App"

function IndexPopup() {
  return (
    <MemoryRouter>
      <Routing />
    </MemoryRouter>
  )
}

export default IndexPopup

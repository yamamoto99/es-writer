import React from "react"
import { Route, Routes } from "react-router-dom"

import Home from "./home"
import SignIn from "./signIn"
import SignUp from "./signUp"
import CheckEmail from "./checkEmail"
import ResendEmail from "./resendEmail"

export default function Routing(){
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/signin" element={<SignIn />} />
      <Route path="/signUp" element={<SignUp />} />
      <Route path="/checkEmail" element={<CheckEmail />} />
      <Route path="/resendEmail" element={<ResendEmail />} />
    </Routes>
  )
}

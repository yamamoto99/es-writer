import React from "react"
import { Route, Routes } from "react-router-dom"

import CheckEmail from "./checkEmail"
import Generating from "./generating"
import Home from "./home"
import LogOut from "./logOut"
import SignIn from "./signIn"
import SignUp from "./signUp"

export default function Routing() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/signin" element={<SignIn />} />
      <Route path="/signUp" element={<SignUp />} />
      <Route path="/checkEmail" element={<CheckEmail />} />
      <Route path="/generating" element={<Generating />} />
      <Route path="/logOut" element={<LogOut />} />
    </Routes>
  )
}

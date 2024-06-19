import React from "react"
import { BrowserRouter as Router, Route, Routes } from "react-router-dom"
import SignUp from "./components/SignUp"
import SignIn from "./components/SignIn"
import ProfileForm from "./components/ProfileForm"
import Profile from "./components/Profile" // プロフィール表示ページが必要な場合

const App = () => {
  console.log("App component rendered");

  return (
    <Router>
      <Routes>
        <Route path="/signup" element={<SignUp />} />
        <Route path="/signin" element={<SignIn />} />
        <Route path="/profileForm" element={<ProfileForm />} />
        <Route path="/profile" element={<Profile />} /> {/* プロフィール表示ページが必要な場合 */}
      </Routes>
    </Router>
  )
}

export default App

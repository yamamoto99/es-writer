import React, { useState } from "react"
import { useStorage } from "@plasmohq/storage/hook"
import { useNavigate } from "react-router-dom"

import openProfileForm from "./openProfileForm"

const signIn = () => {
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const navigate = useNavigate()

  const [loginState, setLoginState] = useStorage<string>("loginState");

  const handleSignIn = async (event: React.FormEvent) => {
    event.preventDefault()
    console.log("SignIn form submitted")

    const response = await fetch("http://35.167.89.55/signin", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ username, password })
    })

    if (response.ok) {
      console.log("SignIn successful")
      setLoginState("logged-in")
      openProfileForm()
    } else {
      console.error("Sign in failed")
      alert("Sign in failed")
    }
  }

  return (
    <>
    <form onSubmit={handleSignIn}>
      <h2>Sign In</h2>
      <input
        type="text"
        placeholder="Username"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
        required
      />
      <input
        type="password"
        placeholder="Password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        required
      />
      <button type="submit">Sign In</button>
    </form>
    <button onClick={() => {setLoginState("not-logged-in");navigate("/")}}>Back</button>
    </>
  )
}

export default signIn

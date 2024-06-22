import React, { useState } from "react"
import openProfileForm from "./openProfileForm"

const signIn = () => {
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")

  const handleSignIn = async (event: React.FormEvent) => {
    event.preventDefault()
    console.log("SignIn form submitted")

    const response = await fetch("http://localhost:8080/signin", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ username, password })
    })

    if (response.ok) {
      console.log("SignIn successful")
      openProfileForm()
    } else {
      console.error("Sign in failed")
      alert("Sign in failed")
    }
  }

  return (
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
  )
}

export default signIn

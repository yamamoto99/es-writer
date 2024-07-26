import React, { useState } from "react"
import { useNavigate } from "react-router-dom"

import { useStorage } from "@plasmohq/storage/hook"

import { api_endpoint } from "../../contents/index"
import openProfileForm from "./openProfileForm"

const signIn = () => {
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const navigate = useNavigate()

  const [loginState, setLoginState] = useStorage<string>("loginState")

  const handleSignIn = async (event: React.FormEvent) => {
    event.preventDefault()
    console.log("SignIn form submitted")

    const response = await fetch(api_endpoint + "/auth/login", {
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
    <form
      onSubmit={handleSignIn}
      className="flex flex-col space-y-1.5 w-40 items-center mb-2 mt-2">
      <input
        type="text"
        placeholder="Username"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
        required
        className="border border-gray-300 rounded-md px-4 py-1 w-5/6"
      />
      <input
        type="password"
        placeholder="Password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        required
        className="border border-gray-300 rounded-md px-4 py-1 w-5/6"
      />
      <div className="flex justify-center space-x-4">
        <button
          type="submit"
          className="bg-blue-500 text-white rounded-md px-3.5 py-2 hover:bg-blue-700">
          Sign In
        </button>
        <button
          onClick={() => {
            setLoginState("not-logged-in")
            navigate("/")
          }}
          className="bg-gray-500 text-white rounded-md px-3 py-2 hover:bg-gray-700">
          Back
        </button>
      </div>
    </form>
  )
}

export default signIn

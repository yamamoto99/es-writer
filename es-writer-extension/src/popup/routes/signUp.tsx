import React, { useState } from "react"
import { useNavigate } from "react-router-dom"

import { useStorage } from "@plasmohq/storage/hook"

import { api_endpoint } from "../../contents/index"

const signUp = () => {
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [email, setEmail] = useState("")
  const navigate = useNavigate()

  const [loginState, setLoginState] = useStorage<string>("loginState")

  const handleSignUp = async (event: React.FormEvent) => {
    event.preventDefault()
    console.log("SignUp form submitted")

    const response = await fetch(api_endpoint + "/auth/signup", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ username, password, email })
    })

    if (response.ok) {
      console.log("SignUp successful")
      setLoginState("checkEmail")
      navigate("/checkEmail")
    } else {
      console.error("Sign up failed")
      alert("Sign up failed")
    }
  }

  return (
    <form
      onSubmit={handleSignUp}
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
      <div className="w-full text-gray-500 text-left px-4">
        <h5>パスワードルール</h5>
        <p>・1つの数字を含む</p>
        <p>・1つの特殊文字を含む</p>
        <p>・1つの大文字を含む</p>
        <p>・1つの小文字を含む</p>
        <p>・8文字以上である</p>
      </div>
      <input
        type="email"
        placeholder="Email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        required
        className="border border-gray-300 rounded-md px-4 py-1 w-5/6"
      />
      <div className="flex justify-center space-x-4">
        <button
          type="submit"
          className="bg-blue-500 text-white rounded-md px-3.5 py-2 hover:bg-blue-700">
          Sign Up
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

export default signUp

import React from "react"
import { useNavigate } from "react-router-dom"

import { useStorage } from "@plasmohq/storage/hook"

import { api_endpoint } from "../../contents"

function LogOut() {
  const navigate = useNavigate()
  const [_, setLoginState] = useStorage<string>("loginState")

  const handleLogout = async () => {
    try {
      const response = await fetch(api_endpoint + "/auth/logout", {
        method: "POST",
        credentials: "include"
      })
      if (response.ok) {
        await setLoginState("not-logged-in")
        navigate("/")
      } else {
        console.error("Logout failed")
      }
    } catch (error) {
      console.error("An error occurred during logout", error)
    }
  }

  return (
    <button
      className="block mx-auto bg-gray-500 hover:bg-gray-700 text-white rounded-md w-24 h-8 p-2 mt-1 mb-4"
      onClick={handleLogout}>
      ログアウト
    </button>
  )
}

export default LogOut

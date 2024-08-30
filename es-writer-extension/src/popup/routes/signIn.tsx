import React from "react"
import { useForm } from "react-hook-form"
import { useNavigate } from "react-router-dom"

import { useStorage } from "@plasmohq/storage/hook"

import { api_endpoint } from "../../contents/index"
import openProfileForm from "./openProfileForm"

const SignIn = () => {
  const navigate = useNavigate()
  const [_, setLoginState] = useStorage<string>("loginState")
  const {
    register,
    handleSubmit,
    formState: { errors }
  } = useForm()

  const onSubmit = async (data) => {
    console.log("SignIn form submitted")
    const response = await fetch(api_endpoint + "/auth/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(data)
    })

    if (response.ok) {
      console.log("SignIn successful")
      setLoginState("logged-in")
      openProfileForm()
    } else {
      console.error("Sign in failed")
      alert("ユーザーネームまたはパスワードが違います")//ここ直す　ユーザー名かパスワードが違う
    }
  }

  return (
    <form
      onSubmit={handleSubmit(onSubmit)}
      className="flex flex-col space-y-1.5 w-40 items-center mb-2 mt-2">
      <input
        type="text"
        placeholder="Username"
        {...register("username", {
          required: "ユーザー名は必須です"
        })}
        className="border border-gray-300 rounded-md px-4 py-1 w-5/6"
      />
      {errors.username && typeof errors.username.message === "string" && (
        <span className="text-red-500 text-xs">{errors.username.message}</span>
      )}

      <input
        type="password"
        placeholder="Password"
        {...register("password", {
          required: "パスワードは必須です"
        })}
        className="border border-gray-300 rounded-md px-4 py-1 w-5/6"
      />
      {errors.password && typeof errors.password.message === "string" && (
        <span className="text-red-500 text-xs">{errors.password.message}</span>
      )}

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

export default SignIn

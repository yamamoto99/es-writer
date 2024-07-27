import React, { useEffect, useState } from "react"
import { useForm } from "react-hook-form"
import { useNavigate } from "react-router-dom"

import { useStorage } from "@plasmohq/storage/hook"

import { api_endpoint } from "../../contents/index"

const SignUp = () => {
  const navigate = useNavigate()
  const [, setLoginState] = useStorage<string>("loginState")
  const [password, setPassword] = useState("")
  const {
    register,
    handleSubmit,
    formState: { errors },
    watch
  } = useForm()

  // パスワードの入力を監視
  useEffect(() => {
    const subscription = watch((value, { name }) => {
      if (name === "password") {
        setPassword(value.password || "")
      }
    })
    return () => subscription.unsubscribe()
  }, [watch])

  // パスワードルールのチェック関数
  const checkPasswordRules = (pass) => ({
    hasNumber: /\d/.test(pass),
    hasSpecialChar: /[\^$*.[\]{}()?"!@#%&\/\\,><':;|_~`=+\-]/.test(pass),
    hasUpperCase: /[A-Z]/.test(pass),
    hasLowerCase: /[a-z]/.test(pass),
    isLongEnough: pass.length >= 8,
    isAsciiPrintable: /^[\x20-\x7E]+$/.test(pass)
  })

  const passwordRules = checkPasswordRules(password)

  // 条件に応じて文字色を返す関数
  const getColorClass = (condition: boolean) =>
    condition ? "text-green-500" : "text-red-500"

  const onSubmit = async (data) => {
    console.log("SignUp form submitted")
    const response = await fetch(api_endpoint + "/auth/signup", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(data)
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
      onSubmit={handleSubmit(onSubmit)}
      className="flex flex-col space-y-1.5 w-60 h-auto items-center mb-2 mt-2">
      <input
        {...register("username", {
          required: "ユーザー名は必須です",
          pattern: {
            value: /^[\x20-\x7E]+$/,
            message: "使用可能文字は半角英数字・記号のみです"
          },
          maxLength: {
            value: 20,
            message: "ユーザー名は20文字以下にしてください"
          }
        })}
        type="text"
        placeholder="Username"
        className="border border-gray-300 rounded-md px-4 py-1 w-5/6"
      />
      {errors.username && typeof errors.username.message === "string" && (
        <span className="text-red-500 text-xs">{errors.username.message}</span>
      )}

      <input
        {...register("email", {
          required: "メールアドレスは必須です",
          pattern: {
            value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
            message: "正しいメールアドレスを入力してください"
          }
        })}
        type="text"
        placeholder="Email"
        className="border border-gray-300 rounded-md px-4 py-1 w-5/6"
      />
      {errors.email && typeof errors.email.message === "string" && (
        <span className="text-red-500 text-xs">{errors.email.message}</span>
      )}

      <input
        {...register("password", {
          required: "パスワードは必須です",
          pattern: {
            value:
              /^(?=.*\d)(?=.*[\^$*.[\]{}()?"!@#%&\/\\,><':;|_~`=+\-])(?=.*[a-z])(?=.*[A-Z])[\x20-\x7E]{8,}$/,
            message: "パスワードが不正です"
          }
        })}
        type="password"
        placeholder="Password"
        className="border border-gray-300 rounded-md px-4 py-1 w-5/6"
      />
      {errors.password && typeof errors.password.message === "string" && (
        <span className="text-red-500 text-xs">{errors.password.message}</span>
      )}
      <div className="w-full text-left px-4">
        <h5 className="text-gray-500">パスワードルール</h5>
        <p className={getColorClass(passwordRules.hasNumber)}>
          {passwordRules.hasNumber ? "○" : "×"} 1つの数字を含む
        </p>
        <p className={getColorClass(passwordRules.hasSpecialChar)}>
          {passwordRules.hasSpecialChar ? "○" : "×"} 1つの特殊文字を含む
        </p>
        <p className={getColorClass(passwordRules.hasUpperCase)}>
          {passwordRules.hasUpperCase ? "○" : "×"} 1つの大文字を含む
        </p>
        <p className={getColorClass(passwordRules.hasLowerCase)}>
          {passwordRules.hasLowerCase ? "○" : "×"} 1つの小文字を含む
        </p>
        <p className={getColorClass(passwordRules.isLongEnough)}>
          {passwordRules.isLongEnough ? "○" : "×"} 8文字以上である
        </p>
        <p className={getColorClass(passwordRules.isAsciiPrintable)}>
          {passwordRules.isAsciiPrintable ? "○" : "×"} 英数字・記号のみである
        </p>
      </div>

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

export default SignUp

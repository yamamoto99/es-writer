import React, { useState } from "react"
import { useNavigate } from "react-router-dom"

import { api_endpoint } from "../../contents/index"

const checkEmail = () => {
  const [verificationCode, setVerificationCode] = useState("")
  const navigate = useNavigate()

  const handleCheckEmail = async (event: React.FormEvent) => {
	  event.preventDefault()
	  console.log("Check Email form submitted")

	  const response = await fetch(api_endpoint + "/auth/checkEmail", {
	    method: "POST",
	    headers: {
	  	"Content-Type": "application/json"
	    },
	    body: JSON.stringify({ verificationCode })
	  })

	  if (response.ok) {
	    console.log("Check Email successful")
	    navigate("/signin")
	  } else {
	    console.error("Check Email failed")
	    alert("Check Email failed")
	  }
  }

  function handleResendEmail() {
	  fetch(api_endpoint + "/auth/resendEmail", {
	    method: "POST",
	  }).then(response => {
      if (response.ok) {
        console.log("Check Email successful")
        alert("Resend Email successful")
      } else {
        console.error("Check Email failed")
        alert("Resend Email failed")
      }
    })
  }

  return (
    <form onSubmit={handleCheckEmail} className="flex flex-col space-y-1.5 w-40 items-center mb-2 mt-2">
      <input
        type="text"
        placeholder="VerificationCode"
        value={verificationCode}
        onChange={(e) => setVerificationCode(e.target.value)}
        required
        className="border border-gray-300 rounded-md px-4 py-1 w-5/6"
      />
      <div className="flex justify-center space-x-3">
        <button
          type="submit"
          className="bg-blue-500 text-white rounded-md px-3 py-2 hover:bg-blue-700"
        >
          Check
        </button>
        <button
          onClick={handleResendEmail}
          type="button"
          className="bg-gray-500 text-white rounded-md px-3 py-2 hover:bg-gray-700"
        >
          resend
        </button>
      </div>
    </form>
  )
}

export default checkEmail

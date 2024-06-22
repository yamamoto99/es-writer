import React, { useState } from "react"
import { useNavigate } from "react-router-dom"

const resendEmail = () => {
  const [verificationCode, setVerificationCode] = useState("")
  const navigate = useNavigate()

  const handleResendEmail = async (event: React.FormEvent) => {
	  event.preventDefault()
	  console.log("Resend Email form submitted")

	  const response = await fetch("http://localhost:8080/checkEmail", {
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

  return (
    <form onSubmit={handleResendEmail}>
      <h2>Check Email</h2>
      <input
        type="text"
        placeholder="VerificationCode"
        value={verificationCode}
        onChange={(e) => setVerificationCode(e.target.value)}
        required
      />
      <button type="submit">Check Email</button>
    </form>
  )
}

export default resendEmail

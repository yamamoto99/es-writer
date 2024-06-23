import React, { useState } from "react"
import { useNavigate } from "react-router-dom"

const checkEmail = () => {
  const [verificationCode, setVerificationCode] = useState("")
  const navigate = useNavigate()

  const handleCheckEmail = async (event: React.FormEvent) => {
	  event.preventDefault()
	  console.log("Check Email form submitted")

	  const response = await fetch("http://35.167.89.55/checkEmail", {
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
	  fetch("http://35.167.89.55/resendEmail", {
	    method: "POST",
	  }).then(response => {
      if (response.ok) {
        console.log("Check Email successful")
      } else {
        console.error("Check Email failed")
        alert("Resend Email failed")
      }
    })
  }

  return (
    <>
    <form onSubmit={handleCheckEmail}>
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
    <button onClick={handleResendEmail}>resend Email</button>
    </>
  )
}

export default checkEmail

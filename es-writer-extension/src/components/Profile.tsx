import React, { useEffect, useState } from "react"

import { api_endpoint } from "../contents/index"

const Profile = () => {
  const [profile, setProfile] = useState({
    bio: "",
    experience: "",
    projects: ""
  })

  useEffect(() => {
    const fetchProfile = async () => {
      const response = await fetch(api_endpoint + "/app/profile/getProfile", {
        method: "GET",
        headers: {
          "Content-Type": "application/json"
        }
      })

      if (response.ok) {
        const data = await response.json()
        setProfile(data)
      } else {
        alert("Failed to fetch profile")
      }
    }

    fetchProfile()
  }, [])

  return (
    <div>
      <h2>Profile</h2>
      <p><strong>Bio:</strong> {profile.bio}</p>
      <p><strong>Experience:</strong> {profile.experience}</p>
      <p><strong>Projects:</strong> {profile.projects}</p>
    </div>
  )
}

export default Profile

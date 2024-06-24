import React, { useState } from "react";

const ProfileForm = () => {
  const [bio, setBio] = useState("");
  const [experience, setExperience] = useState("");
  const [projects, setProjects] = useState("");

  const handleProfileSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    const response = await fetch(api_endpoint + "/saveprofile", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ bio, experience, projects })
    });

    if (response.ok) {
      alert("Profile saved successfully");
    } else {
      alert("Failed to save profile");
    }
  };

  const formStyle = {
    maxWidth: "600px",
    margin: "0 auto"
  };

  const labelStyle = {
    display: "block",
    marginBottom: "10px"
  };

  const textareaStyle = {
    width: "100%",
    height: "100px",
    marginBottom: "20px",
  };

  const buttonStyle = {
    display: "block",
    margin: "0 auto",
  };

  return (
    <form onSubmit={handleProfileSubmit} style={formStyle}>
      <h2>Profile Information</h2>
      <label style={labelStyle}>
        自己PR:
        <textarea
          value={bio}
          onChange={(e) => setBio(e.target.value)}
          required
          style={textareaStyle}
        />
      </label>
      <label style={labelStyle}>
        経験:
        <textarea
          value={experience}
          onChange={(e) => setExperience(e.target.value)}
          required
          style={textareaStyle}
        />
      </label>
      <label style={labelStyle}>
        今まで作った作品:
        <textarea
          value={projects}
          onChange={(e) => setProjects(e.target.value)}
          required
          style={textareaStyle}
        />
      </label>
      <button type="submit" style={buttonStyle}>Save Profile</button>
    </form>
  );
};

export default ProfileForm;

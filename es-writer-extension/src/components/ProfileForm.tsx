import React, { useState } from "react";

const ProfileForm = () => {
  const [bio, setBio] = useState("");
  const [experience, setExperience] = useState("");
  const [projects, setProjects] = useState("");

  const handleProfileSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    const response = await fetch("http://localhost:8080/saveProfile", {
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

  return (
    <form onSubmit={handleProfileSubmit}>
      <h2>Profile Information</h2>
      <label>
        Bio:
        <textarea
          value={bio}
          onChange={(e) => setBio(e.target.value)}
          required
        />
      </label>
      <label>
        Experience:
        <textarea
          value={experience}
          onChange={(e) => setExperience(e.target.value)}
          required
        />
      </label>
      <label>
        Projects:
        <textarea
          value={projects}
          onChange={(e) => setProjects(e.target.value)}
          required
        />
      </label>
      <button type="submit">Save Profile</button>
    </form>
  );
};

export default ProfileForm;

import React, { useState, useEffect } from "react";

import "../../style.css";

import { api_endpoint } from "../contents/index"

const ProfileForm = () => {
  const [bio, setBio] = useState("");
  const [experience, setExperience] = useState("");
  const [projects, setProjects] = useState("");

  useEffect(() => {
    const fetchProfileData = async () => {
      try {
        const response = await fetch(api_endpoint + "/app/profile/getProfile", {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
          }
        });

        if (response.ok) {
          const data = await response.json();
          setBio(data.bio || "");
          setExperience(data.experience || "");
          setProjects(data.projects || "");
        } else {
          console.error("Failed to fetch profile data");
        }
      } catch (error) {
        console.error("Error fetching profile data:", error);
      }
    };

    fetchProfileData();
  }, []);

  const handleProfileSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    const response = await fetch(api_endpoint + "/app/profile/updateProfile", {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
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
    <form onSubmit={handleProfileSubmit} className="max-w-lg mx-auto">
      <h2 className="text-xl font-bold mb-4">Profile Information</h2>
      <div className="mb-4">
        <label className="block mb-2">
          自己PR:
          <textarea
            value={bio}
            onChange={(e) => setBio(e.target.value)}
            required
            className="w-full h-24 p-2 border border-gray-300 rounded mb-4"
          />
        </label>
      </div>
      <div className="mb-4">
        <label className="block mb-2">
          経験:
          <textarea
            value={experience}
            onChange={(e) => setExperience(e.target.value)}
            required
            className="w-full h-24 p-2 border border-gray-300 rounded mb-4"
          />
        </label>
      </div>
      <div className="mb-4">
        <label className="block mb-2">
          今まで作った作品:
          <textarea
            value={projects}
            onChange={(e) => setProjects(e.target.value)}
            required
            className="w-full h-24 p-2 border border-gray-300 rounded mb-4"
          />
        </label>
      </div>
      <button type="submit" className="block mx-auto px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-700">Save Profile</button>
    </form>
  );
};

export default ProfileForm;

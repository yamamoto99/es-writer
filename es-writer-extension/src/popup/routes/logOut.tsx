import React from "react";
import { useStorage } from "@plasmohq/storage/hook"
import { useNavigate } from "react-router-dom";
import { api_endpoint } from "../../contents";

function LogOut() {
  const navigate = useNavigate();
  const [loginState, setLoginState] = useStorage<string>("loginState");

  const handleLogout = async () => {
    try {
      const response = await fetch(api_endpoint + "/auth/logout", {
        method: "POST",
        credentials: "include", // Cookie を含める
      });
      if (response.ok) {
        setLoginState("not-logged-in");
        navigate("/");
      } else {
        console.error("Logout failed");
      }
    } catch (error) {
      console.error("An error occurred during logout", error);
    }
  };

  return (
    <button onClick={handleLogout}>
      Logout
    </button>
  );
}

export default LogOut;

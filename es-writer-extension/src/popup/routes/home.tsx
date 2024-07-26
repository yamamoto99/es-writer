import React, { useEffect } from "react";

import { useNavigate } from "react-router-dom";
import { useStorage } from "@plasmohq/storage/hook"

import genAnswer from "./genAnswer";
import openProfileForm from "./openProfileForm";
import { api_endpoint } from "../../contents/index"
import "../../../style.css";

async function fetchData(loginState: string | undefined, setLoginState: (loginState: string) => void){
  try {
    const response = await fetch(api_endpoint + "/auth/login", {
      method: "POST"
    });

    if (response.ok) {
      setLoginState("logged-in");
    } else {
      if (typeof loginState === "undefined" || loginState === "not-logged-in" || loginState === "logged-in") {
        setLoginState("not-logged-in");
      }
    }
  } catch (error) {
    console.error("Fetch error:", error);
  }
}

function IndexPopup() {
  const navigate = useNavigate();

  const [loginState, setLoginState] = useStorage<string>("loginState");

  useEffect(() => {
    fetchData(loginState, setLoginState);
  }, []);

  if (loginState === "not-logged-in") {
    return (
      <div className="w-40 h-20">
        <button 
          className="block mx-auto bg-blue-500 hover:bg-blue-700 text-white rounded-md w-32 h-8 p-2 mt-4 mb-1"
          onClick={() => {setLoginState("signUp");navigate("/signUp")}}
        >
          サインアップ
        </button>
        <button
          className="block mx-auto bg-blue-500 hover:bg-blue-700 text-white rounded-md w-32 h-8 p-2 mt-1 mb-4"
          onClick={() => {setLoginState("signIn");navigate("/signIn")}}
        >
          サインイン
        </button>
      </div>
    );
  } else if (loginState === "logged-in") {
    return (
      <div className="w-40 h-20">
        <button
          className="block mx-auto bg-blue-500 hover:bg-blue-700 text-white rounded-md w-32 h-8 p-2 mt-4 mb-1"
          onClick={async () => {
            navigate("/generating");
            try {
              await genAnswer();
            } catch (error) {
              console.error("Error generating answer:", error);
            }
            window.close();
          }}
        >
          回答生成
        </button>
        <button
          className="block mx-auto bg-blue-500 hover:bg-blue-700 text-white rounded-md w-32 h-8 p-2 mt-1 mb-4"
          onClick={openProfileForm}
        >
          経歴入力
        </button>
        <button
          className="block mx-auto bg-gray-500 hover:bg-gray-700 text-white rounded-md w-24 h-6 p-1 mt-1 mb-2"
          onClick={() => { setLoginState("not-logged-in"); navigate("/logOut") }}
        >
          ログアウト
        </button>
      </div>
    );
  }else if (loginState === "signUp") {
    navigate("/signUp");
  }else if (loginState === "signIn") {
    navigate("/signIn");
  }else if (loginState === "checkEmail") {
    navigate("/checkEmail");
  }

}

export default IndexPopup;

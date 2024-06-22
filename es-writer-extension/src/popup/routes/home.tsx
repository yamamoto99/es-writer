import React from "react";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

import genAnswer from "./genAnswer";
import openProfileForm from "./openProfileForm";

function IndexPopup() {
  const navigation = useNavigate();

  // /welcomeをたたいて、成功すれば、ログイン済みとして扱う
  const [isLogin, setIsLogin] = useState(false);
  fetch("http://localhost:8080/welcome", {
    method: "GET"
  }).then((response) => {
    if (response.ok) {
      setIsLogin(true);
    } else {
      setIsLogin(false);
    }
  });

  if (!isLogin) {
    return (
      <div style={{ width: '150px', height: '75px' }}>
        <button onClick={() => navigation("/signUp")}>サインアップ</button>
        <button onClick={() => navigation("/signIn")}>サインイン</button>
      </div>
    );
  }else{
    return (
      <div style={{ width: '150px', height: '75px' }}>
        <button onClick={genAnswer}>回答生成</button>
        <button onClick={openProfileForm}>経歴入力</button>
      </div>
    );
  }
}

export default IndexPopup;

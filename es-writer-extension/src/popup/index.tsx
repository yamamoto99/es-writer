import React from "react";
import { useState } from "react";
import indexContents from "./indexContents";
import openProfileForm from "./openProfileForm";

function IndexPopup() {
  // /welcomeをたたいて、成功すれば、ログイン済みとして扱う
  // 失敗した場合、ログインしていないとして扱い、サインアップまたはサインインボタンを表示する
  const [isLogin, setIsLogin] = useState(false);
  // /welcomeをたたいてレスポンスを受け取る
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
        <button onClick={openProfileForm}>サインアップ</button>
        <button onClick={openProfileForm}>サインイン</button>
      </div>
    );
  }else{
    return (
      <div style={{ width: '150px', height: '75px' }}>
        <button onClick={indexContents}>回答生成</button>
        <button onClick={openProfileForm}>経歴入力</button>
      </div>
    );
  }
}

export default IndexPopup;

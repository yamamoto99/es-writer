## 拡張機能の動かし方

ます、以下を実行

```
yarn install
```

次にビルド

```
yarn build
```

Chromeで開発者モードをオンにした上で、`es-writer-extension/build/chrome-XXX-XXXX`ディレクトリの拡張機能を読み込ませる。

右上の拡張機能から`回答生成`でWebAPIへ取得しに行き、textArea内の入力文字を置き換える(動作未確認)

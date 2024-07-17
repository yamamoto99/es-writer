## 拡張機能の動かし方

pnpmをインストールする

```
npm install -g pnpm
```

ます、以下を実行

```
pnpm install
```

次にビルド

```
pnpm build
```

Chromeで開発者モードをオンにした上で、`es-writer-extension/build/chrome-XXX-XXXX`ディレクトリの拡張機能を読み込ませる。

サインアップ or サインインを行い、`回答生成`を押すと、textArea内に自動で入力される。

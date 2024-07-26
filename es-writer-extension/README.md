## 拡張機能の動かし方

pnpmをインストールする

```
npm install -g pnpm
```

ます、以下を実行

```
pnpm install
```

次にデベロップサーバーを立ち上げる

```
pnpm dev
```

Chromeで開発者モードをオンにした上で、`es-writer-extension/build/chrome-mv3-dev`ディレクトリの拡張機能を読み込ませる。

サインアップ or サインインを行い、`回答生成`を押すと、textArea内に自動で入力される。

## コマンド＆オプション一覧

- プロダクションビルドを行う

```
pnpm build
```

- デベロップサーバーをたてる（HMRがきく）

```
pnpm dev
```

- プロダクションビルドをzipファイルにする

```
pnpm build
pnpm package
```

or

```
pnpm build --zip
```

- ターゲットのブラウザを指定するオプション（デフォルトはchromeのmanifest v3）

```
--target=firefox-mv3
```

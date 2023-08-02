# 概要

morayは、AWS System Mnagerセッションマネージャーでリモートホストへのポートフォワードをする操作を楽にするために作成したCLIツールです。

フォワーディング先の対象は`RDS`及び`DocumentDB`となります。

# 前提条件

## 依存ツール

morayを使用するためには、AWSの`Session Manager Plugin`がインストールされている必要があります。

[AWS CLI 用の Session Manager プラグインをインストールする](https://docs.aws.amazon.com/ja_jp/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)

## AWS IAM権限

morayを実行するにあたり、以下のAWS IAMポリシーが必要です。

```json
工事中
```

# インストール

## Homebrew (macOS and Linux)

```bash
brew install tomozo6/tap/moray
```

## Scoop (Windows PowerShell)

```bash
scoop bucket add tomozo6 https://github.com/tomozo6/scoop-bucket
scoop install moray
```

## Binary Packages

Download from [Releases](https://github.com/tomozo6/moray/releases).

# Usage

よく使用されるコマンドをいくつか紹介します。詳細なフラグについては`moray [command] --help`を実行して確認してください。

## moray

```bash
moray
```

踏み台EC2や接続先のDBを対話形式で選択し、ポートフォワーディングをします。
接続先のDBは`リーダーインスタンス`になります。

なお接続先DBのポート番号と同じ番号をローカルポートとしてフォワーディングします。

```bash
moray --witer
```

`--writer`フラグを使用すると、接続先のDBは`ライターインスタンス`になります。

```bash
moray --profile stg
```

`--profile`フラグを使用すると、自身で設定した任意の名前付きプロファイルを参照します。

```bash
moray --port 3333
```

`--port`フラグを使用すると、任意のローカルポートを使用してフォワーディングします。

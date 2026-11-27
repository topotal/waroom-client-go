# Security Policy

[English](#english) | [日本語](#日本語)

## English

### Supported Versions

This project is in early development (v0.x). Security fixes are provided for the latest release only.
The first public release is v0.1.0. Earlier versions (v0.0.x) were released before the repository was made public and are not supported.

| Version                          | Supported          |
| -------------------------------- | ------------------ |
| Latest release (v0.1.0 or later) | :white_check_mark: |
| Older releases                   | :x:                |

### Reporting a Vulnerability

Report vulnerabilities privately through GitHub Security Advisories
([Private vulnerability reporting](https://docs.github.com/en/code-security/security-advisories/guidance-on-reporting-and-writing-information-about-vulnerabilities/privately-reporting-a-security-vulnerability)):

1. Open the [Security tab](https://github.com/topotal/waroom-client-go/security) of this repository.
2. Click **Report a vulnerability** and fill in the form.

> [!IMPORTANT]
> Do not report vulnerabilities through public issues, pull requests, or discussions.
> They are visible to everyone and may put users at risk before a fix is available.

Please include as much of the following as you can:

- Affected version(s)
- Description of the vulnerability and its impact
- Steps to reproduce or a proof of concept
- Suggested fix, if any

Do not include real API keys or data belonging to any Waroom organization. Use placeholders or redacted values instead.

Reports may be written in English or Japanese.

### Scope

This policy covers waroom-client-go. If you are not sure whether an issue is in this library or in the Waroom service itself, report it here and we will forward it to the appropriate team.

### What to Expect

We will acknowledge your report within 5 business days (Japan time).

After confirming the vulnerability, we will work on a fix and coordinate the disclosure timing with you.
We publish a GitHub Security Advisory after a fixed version is released. Please do not disclose the vulnerability publicly until the advisory is published.

If you wish, we will credit you as the reporter in the advisory.

## 日本語

### サポート対象バージョン

本プロジェクトは開発初期段階 (v0.x) です。セキュリティ修正は最新リリースに対してのみ提供します。
最初の公開リリースは v0.1.0 です。それ以前のバージョン (v0.0.x) はリポジトリ公開前のリリースであり、サポート対象外です。

| バージョン                     | サポート           |
| ------------------------------ | ------------------ |
| 最新リリース (v0.1.0 以降)     | :white_check_mark: |
| それ以前のリリース             | :x:                |

### 脆弱性の報告

脆弱性は GitHub Security Advisories
([Private vulnerability reporting](https://docs.github.com/ja/code-security/security-advisories/guidance-on-reporting-and-writing-information-about-vulnerabilities/privately-reporting-a-security-vulnerability))
から非公開で報告してください。

1. 本リポジトリの [Security タブ](https://github.com/topotal/waroom-client-go/security) を開きます。
2. **Report a vulnerability** をクリックし、フォームに記入します。

> [!IMPORTANT]
> 公開の Issue、Pull Request、Discussion では脆弱性を報告しないでください。
> 誰でも閲覧できるため、修正が提供される前に利用者を危険にさらすおそれがあります。

可能な範囲で以下の情報を含めてください。

- 影響を受けるバージョン
- 脆弱性の内容と影響
- 再現手順または PoC
- 修正案 (あれば)

実際の API キーや、Waroom の組織に属するデータは含めないでください。プレースホルダーやマスクした値を使用してください。

報告は日本語・英語のどちらでも受け付けます。

### 対象範囲

本ポリシーは waroom-client-go を対象とします。問題が本ライブラリと Waroom サービス本体のどちらにあるか分からない場合も、こちらから報告してください。適切なチームに引き継ぎます。

### 報告後の対応

報告の受領後、5 営業日 (日本時間) 以内に一次返信します。

脆弱性を確認した後、修正に取り組み、公開時期は報告者と相談して決めます。
修正版をリリースした後に GitHub Security Advisory を公開します。Advisory が公開されるまでは、脆弱性を公開しないようお願いします。

希望される場合は、Advisory に報告者としてクレジットを記載します。

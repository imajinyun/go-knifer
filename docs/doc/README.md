# 📚 go-knifer Documentation Hub / 文档中心

> Detailed package navigation, architecture notes, safety defaults, and contribution workflow for `go-knifer`.<br>
> `go-knifer` 的详细包导航、架构说明、安全默认值与贡献工作流。

## 📑 Table of Contents / 目录

- [🧭 Quick navigation / 快速导航](#quick-navigation)
- [🧩 Package catalog / 包目录](#package-catalog)
- [📝 Quickstart documents / 快速开始文档](#quickstart-documents)
- [🏗️ Architecture and package boundaries / 架构与包边界](#architecture-and-package-boundaries)
- [✅ Recommended API entry points / 推荐 API 入口](#recommended-api-entry-points)
- [📦 Build, test, and release workflow / 构建、测试与发布工作流](#build-test-and-release-workflow)
- [🛡️ Governance / 治理](#governance)
- [🤝 Contributing / 贡献](#contributing)

<a id="quick-navigation"></a>

## 🧭 Quick navigation / 快速导航

- 🏠 Root README / 根 README：[`../../README.md`](../../README.md)
- 🌐 Online Go docs / 在线 Go 文档：[pkg.go.dev/github.com/imajinyun/go-knifer](https://pkg.go.dev/github.com/imajinyun/go-knifer)
- 🧾 Public API snapshot / 公开 API 快照：[`../api/exports.txt`](../api/exports.txt)
- 🗺️ AI-oriented project map / 面向 AI 的项目地图：[`../../llms.txt`](../../llms.txt)
- 🧯 Security policy / 安全策略：[`../../SECURITY.md`](../../SECURITY.md)
- 📝 Changelog / 变更日志：[`../../CHANGELOG.md`](../../CHANGELOG.md)

<a id="package-catalog"></a>

## 🧩 Package catalog / 包目录

The project follows an “internal implementation + public facade” layout: `internal/*` contains concrete implementations, while `v*` packages expose stable public APIs.

项目采用“内部实现 + 公开 facade”的结构：`internal/*` 保存具体实现，`v*` 包提供稳定公开 API。

| Module / 模块 | Import path / 导入路径 | Description / 说明 |
| --- | --- | --- |
| [`vbean`](01-vbean.md) | `github.com/imajinyun/go-knifer/vbean` | Bean and struct mapping helpers.<br>Bean 与结构体映射工具。 |
| [`vblf`](02-vblf.md) | `github.com/imajinyun/go-knifer/vblf` | Bloom filter helpers.<br>布隆过滤器工具。 |
| [`vbool`](03-vbool.md) | `github.com/imajinyun/go-knifer/vbool` | Boolean helpers.<br>布尔工具。 |
| [`vcache`](04-vcache.md) | `github.com/imajinyun/go-knifer/vcache` | FIFO/LFU/LRU/timed/weak cache helpers.<br>FIFO/LFU/LRU/定时/弱引用缓存工具。 |
| [`vcodec`](05-vcodec.md) | `github.com/imajinyun/go-knifer/vcodec` | Base64 and Hex encoding helpers.<br>Base64 与 Hex 编解码工具。 |
| [`vconf`](06-vconf.md) | `github.com/imajinyun/go-knifer/vconf` | Configuration loading and validation helpers.<br>配置加载与校验工具。 |
| [`vconv`](07-vconv.md) | `github.com/imajinyun/go-knifer/vconv` | Permissive type conversion helpers.<br>宽松类型转换工具。 |
| [`vcron`](08-vcron.md) | `github.com/imajinyun/go-knifer/vcron` | Cron parsing and scheduling helpers.<br>Cron 解析与调度工具。 |
| [`vcrypto`](09-vcrypto.md) | `github.com/imajinyun/go-knifer/vcrypto` | Cryptography, digest, signing, and key helpers.<br>密码学、摘要、签名与密钥工具。 |
| [`vcsv`](10-vcsv.md) | `github.com/imajinyun/go-knifer/vcsv` | CSV reading and writing helpers.<br>CSV 读写工具。 |
| [`vdate`](11-vdate.md) | `github.com/imajinyun/go-knifer/vdate` | Date and time helpers.<br>日期时间工具。 |
| [`vdb`](12-vdb.md) | `github.com/imajinyun/go-knifer/vdb` | Database/sql helper APIs.<br>database/sql 辅助 API。 |
| [`vdfa`](13-vdfa.md) | `github.com/imajinyun/go-knifer/vdfa` | DFA word-tree matching helpers.<br>DFA 词树匹配工具。 |
| [`verr`](14-verr.md) | `github.com/imajinyun/go-knifer/verr` | Error aggregation, recovery, and stack helpers.<br>错误聚合、恢复与堆栈工具。 |
| [`vfile`](15-vfile.md) | `github.com/imajinyun/go-knifer/vfile` | File and IO helpers.<br>文件与 IO 工具。 |
| [`vform`](16-vform.md) | `github.com/imajinyun/go-knifer/vform` | Form and input validation helpers.<br>表单与输入校验工具。 |
| [`vhash`](17-vhash.md) | `github.com/imajinyun/go-knifer/vhash` | Non-cryptographic hash helpers.<br>非密码学哈希工具。 |
| [`vhttp`](18-vhttp.md) | `github.com/imajinyun/go-knifer/vhttp` | Standard-library HTTP facade.<br>标准库 HTTP facade。 |
| [`vid`](19-vid.md) | `github.com/imajinyun/go-knifer/vid` | UUID, Snowflake, ObjectId, and NanoId helpers.<br>UUID、Snowflake、ObjectId 与 NanoId 工具。 |
| [`vident`](20-vident.md) | `github.com/imajinyun/go-knifer/vident` | Identity number parsing and validation helpers.<br>身份号码解析与校验工具。 |
| [`vimg`](21-vimg.md) | `github.com/imajinyun/go-knifer/vimg` | Image, QR/barcode, and captcha helpers.<br>图像、二维码/条码与验证码工具。 |
| [`vjob`](22-vjob.md) | `github.com/imajinyun/go-knifer/vjob` | Sliceable job execution helpers.<br>可切片任务执行工具。 |
| [`vjson`](23-vjson.md) | `github.com/imajinyun/go-knifer/vjson` | JSON object, array, path, and conversion helpers.<br>JSON 对象、数组、路径与转换工具。 |
| [`vjwt`](24-vjwt.md) | `github.com/imajinyun/go-knifer/vjwt` | JWT signing, parsing, and verification helpers.<br>JWT 签名、解析与校验工具。 |
| [`vlog`](25-vlog.md) | `github.com/imajinyun/go-knifer/vlog` | Logging facade helpers.<br>日志 facade 工具。 |
| [`vmail`](26-vmail.md) | `github.com/imajinyun/go-knifer/vmail` | Mail composition and SMTP sending helpers.<br>邮件构建与 SMTP 发送工具。 |
| [`vmap`](27-vmap.md) | `github.com/imajinyun/go-knifer/vmap` | Map construction, transform, merge, and diff helpers.<br>Map 构造、转换、合并与比较工具。 |
| [`vmask`](28-vmask.md) | `github.com/imajinyun/go-knifer/vmask` | Sensitive data masking helpers.<br>敏感数据脱敏工具。 |
| [`vnet`](29-vnet.md) | `github.com/imajinyun/go-knifer/vnet` | Network, IP, CIDR, TLS, and multipart helpers.<br>网络、IP、CIDR、TLS 与 multipart 工具。 |
| [`vnum`](30-vnum.md) | `github.com/imajinyun/go-knifer/vnum` | Numeric, arithmetic, range, and expression helpers.<br>数字、算术、范围与表达式工具。 |
| [`vobj`](31-vobj.md) | `github.com/imajinyun/go-knifer/vobj` | Object-level nil, empty, clone, and compare helpers.<br>对象级 nil、empty、clone 与比较工具。 |
| [`vpass`](32-vpass.md) | `github.com/imajinyun/go-knifer/vpass` | Password strength analysis helpers.<br>密码强度分析工具。 |
| [`vpoi`](33-vpoi.md) | `github.com/imajinyun/go-knifer/vpoi` | Office/XLSX document helpers.<br>Office/XLSX 文档工具。 |
| [`vrand`](34-vrand.md) | `github.com/imajinyun/go-knifer/vrand` | Random value and secure byte helpers.<br>随机值与安全随机字节工具。 |
| [`vref`](35-vref.md) | `github.com/imajinyun/go-knifer/vref` | Reflection field, method, and type helpers.<br>反射字段、方法与类型工具。 |
| [`vregex`](36-vregex.md) | `github.com/imajinyun/go-knifer/vregex` | Regular-expression helper APIs.<br>正则表达式工具 API。 |
| [`vresty`](37-vresty.md) | `github.com/imajinyun/go-knifer/vresty` | Resty-based HTTP facade.<br>基于 Resty 的 HTTP facade。 |
| [`vsem`](38-vsem.md) | `github.com/imajinyun/go-knifer/vsem` | Weighted semaphore helpers.<br>带权信号量工具。 |
| [`vset`](39-vset.md) | `github.com/imajinyun/go-knifer/vset` | Generic and typed set helpers.<br>泛型与类型化集合工具。 |
| [`vskt`](40-vskt.md) | `github.com/imajinyun/go-knifer/vskt` | TCP socket client/server helpers.<br>TCP socket 客户端/服务端工具。 |
| [`vslice`](41-vslice.md) | `github.com/imajinyun/go-knifer/vslice` | Slice query, transform, set-style, and paging helpers.<br>Slice 查询、转换、集合式操作与分页工具。 |
| [`vstr`](42-vstr.md) | `github.com/imajinyun/go-knifer/vstr` | String and text processing helpers.<br>字符串与文本处理工具。 |
| [`vsys`](43-vsys.md) | `github.com/imajinyun/go-knifer/vsys` | System and runtime information helpers.<br>系统与运行时信息工具。 |
| [`vtpl`](44-vtpl.md) | `github.com/imajinyun/go-knifer/vtpl` | html/template rendering helpers.<br>html/template 渲染工具。 |
| [`vurl`](45-vurl.md) | `github.com/imajinyun/go-knifer/vurl` | URL/URI parsing, encoding, and safe-resource helpers.<br>URL/URI 解析、编码与安全资源工具。 |
| [`vver`](46-vver.md) | `github.com/imajinyun/go-knifer/vver` | Version comparison and expression helpers.<br>版本比较与表达式工具。 |
| [`vxml`](47-vxml.md) | `github.com/imajinyun/go-knifer/vxml` | XML parse, tree, format, and conversion helpers.<br>XML 解析、树操作、格式化与转换工具。 |
| [`vzip`](48-vzip.md) | `github.com/imajinyun/go-knifer/vzip` | ZIP/gzip/zlib archive and compression helpers.<br>ZIP/gzip/zlib 归档与压缩工具。 |

<a id="quickstart-documents"></a>

## 📝 Quickstart documents / 快速开始文档

Per-package quickstart examples live in the linked documents above so examples stay focused and easy to maintain by domain.

每个包的 quickstart 示例都放在上方链接的独立文档中，便于按领域维护，也让示例更聚焦。

<a id="architecture-and-package-boundaries"></a>

## 🏗️ Architecture and package boundaries / 架构与包边界

`go-knifer` uses public `v*` packages as facade APIs and keeps concrete code in `internal/*`. Application code should import the `v*` packages; `internal/*` exists so implementations can evolve without exposing every helper as public API.

`go-knifer` 使用公开 `v*` 包作为 facade API，并将具体实现保留在 `internal/*`。业务代码应导入 `v*` 包；`internal/*` 用于让实现可以演进，而不必把每个 helper 都暴露为公开 API。

Facade rules / facade 规则：

- `internal/<domain>` owns implementation details and domain-specific tests. / `internal/<domain>` 负责实现细节和领域测试。
- `v<domain>` exposes the stable public surface for that domain. / `v<domain>` 暴露该领域稳定的公开 API。
- Small utility packages may use hand-written thin facades; larger modules may keep generated `facade.go` files. / 小工具包可以使用手写薄 facade；较大模块可以保留生成的 `facade.go`。
- Newly exported internal APIs should be reviewed before being exposed publicly. / 新增 internal 导出 API 在公开暴露前应先评审。
- Short facade names such as `vform`, `vmask`, `vsem`, `vskt`, `vblf`, and `vver` are documented instead of renamed. / `vform`、`vmask`、`vsem`、`vskt`、`vblf`、`vver` 等短包名通过文档解释，不随意重命名。

API compatibility / API 兼容性：

- The root package and top-level `v*` subpackages are the public compatibility boundary. / 根包与顶层 `v*` 子包是公开兼容边界。
- Their exported API surface is recorded in [`../api/exports.txt`](../api/exports.txt). / 导出 API 面记录在 [`../api/exports.txt`](../api/exports.txt)。
- Standalone `internal/*` packages are intentionally excluded so implementation packages can be refactored without public API noise. / 独立 `internal/*` 包有意排除在外，便于重构实现而不产生公开 API 噪音。
- `make api-check` regenerates a temporary snapshot and compares it with the checked-in file. / `make api-check` 会生成临时快照并与仓库中的快照比较。
- For intentional public API changes, run `UPDATE_API=1 make api-check` and review the snapshot diff. / 如需有意变更公开 API，运行 `UPDATE_API=1 make api-check` 并审查快照 diff。

Configurable APIs and provider injection / 可配置 API 与 provider 注入：

- Many packages expose functional options through `WithXxx` helpers and `XxxWithOptions` variants. / 许多包通过 `WithXxx` helper 和 `XxxWithOptions` 变体暴露函数式选项。
- Existing fixed-argument APIs stay stable, while option-based variants add advanced control. / 既有固定参数 API 保持稳定，基于 option 的变体提供高级控制能力。
- Provider-style options let callers inject dependencies for deterministic tests and controlled runtime behavior. / provider 风格选项允许调用方注入依赖，以支持确定性测试和受控运行时行为。
- Package-level defaults remain explicit and should not be mutated implicitly by per-call options. / 包级默认值保持显式，按调用传入的 option 不应隐式修改全局状态。

Domain boundary rules / 领域边界规则：

- `vhash` owns non-cryptographic hash helpers; `vcrypto` owns security-oriented cryptographic operations. / `vhash` 负责非密码学哈希；`vcrypto` 负责安全相关密码学操作。
- `vhttp` is the standard-library HTTP facade; `vresty` is the Resty-based facade. / `vhttp` 是标准库 HTTP facade；`vresty` 是 Resty facade。
- URL escaping, query building/parsing, and scheme checks live in `vurl`. / URL 转义、查询构造/解析和 scheme 检查归属 `vurl`。
- `vid` owns generated identifiers; `vident` owns legal identity numbers and regional card parsing. / `vid` 负责生成型 ID；`vident` 负责法定身份号码与地区证件解析。
- `vobj` is a convenience object-level facade; new domain logic should still start in focused packages. / `vobj` 是对象级便利 facade，新领域逻辑仍应优先放入聚焦包。

### 🚦 Error contract / 错误契约

The root package `knifer` owns the cross-cutting error contract: the `ErrCode` classifier, unified `knifer.Error` type, `CodeCarrier` interface, `CodeOf` extractor, and `NewError` / `WrapError` / `Errorf` constructors.

根包 `knifer` 负责跨包错误契约：`ErrCode` 分类器、统一 `knifer.Error` 类型、`CodeCarrier` 接口、`CodeOf` 提取器，以及 `NewError` / `WrapError` / `Errorf` 构造器。

```go
if errors.Is(err, knifer.ErrCodeInvalidInput) { /* ... */ }
if code, ok := knifer.CodeOf(err); ok { /* ... */ }
```

`vcrypto` is a reference integration: validation errors match both `knifer.ErrCodeInvalidInput` and existing `vcrypto` sentinels.<br>
`vcrypto` 是参考集成：校验错误既匹配 `knifer.ErrCodeInvalidInput`，也匹配现有 `vcrypto` 哨兵错误。

### 🔐 Security and safety defaults / 安全默认值

Security-sensitive helpers expose only the currently recommended public API surface.<br>
安全敏感 helper 只暴露当前推荐的公开 API 面。

- `vcrypto` keeps SHA-2, HMAC, PBKDF2-SHA256, AES-GCM, RSA-OAEP, and RSA-PSS. / `vcrypto` 保留 SHA-2、HMAC、PBKDF2-SHA256、AES-GCM、RSA-OAEP 和 RSA-PSS。
- JWT RSA signing uses RSA-PSS, and unsigned JWT `alg=none` tokens are rejected. / JWT RSA 签名使用 RSA-PSS，并拒绝 unsigned `alg=none` token。
- TLS helpers use TLS 1.2+ minimum versions and do not bypass certificate verification. / TLS helper 使用 TLS 1.2+ 最低版本，不绕过证书校验。
- Safe HTTP/URL helpers reject local/private/link-local/unspecified targets by default and re-check redirects. / 安全 HTTP/URL helper 默认拒绝本地、私有、link-local、未指定地址，并重新检查重定向。
- `vfile`, `vconf`, `vurl`, and `vzip` use bounded reads or extraction/decompression limits by default. / `vfile`、`vconf`、`vurl`、`vzip` 默认使用有界读取或解压限制。

<a id="recommended-api-entry-points"></a>

## ✅ Recommended API entry points / 推荐 API 入口

Use these APIs for new code. Request helpers that can fail return errors explicitly instead of swallowing failures.

新代码建议使用以下 API。可能失败的请求 helper 应显式返回错误，而不是吞掉失败。

| Scenario / 场景 | Recommended API / 推荐 API |
| --- | --- |
| Build a trusted standard-library HTTP request<br>构建可信标准库 HTTP 请求 | `vhttp.Get`, `vhttp.Post`, `vhttp.NewRequest` |
| Read a trusted HTTP response body and handle errors<br>读取可信 HTTP 响应体并处理错误 | `vhttp.GetStringE`, `vhttp.PostJSONE`, `vhttp.DownloadBytesE` |
| Access user-controlled or otherwise untrusted HTTP(S) URLs<br>访问用户可控或不可信 HTTP(S) URL | `vhttp.GetStringSafeE`, `vhttp.PostJSONSafeE`, `vhttp.DownloadBytesSafeE`, `vurl.OpenSafe` |
| Use the Resty-backed HTTP facade<br>使用 Resty HTTP facade | `vresty.Get`, `vresty.Post`, `vresty.GetStringE`, `vresty.PostJSONE` |
| Access untrusted URLs through Resty<br>通过 Resty 访问不可信 URL | `vresty.GetStringSafeE`, `vresty.PostJSONSafeE`, `vresty.DownloadBytesSafeE` |
| Download a user-controlled URL to a file<br>下载用户可控 URL 到文件 | `vhttp.DownloadFileSafe` or `vresty.DownloadFileSafe` |
| Generate bytes for secrets, tokens, keys, nonces, or salts<br>生成密钥、令牌、key、nonce 或 salt 字节 | `vrand.SecureBytes` |
| Create an LRU cache<br>创建 LRU 缓存 | `vcache.NewLRU` or `vcache.NewLRUWithTimeout` |
| Parse a cron expression<br>解析 cron 表达式 | `vcron.NewPattern` or `vcron.MustNewPattern` |
| Load remote configuration from a trust boundary<br>从信任边界加载远程配置 | `vconf.LoadRemoteSafe` or `vconf.LoadRemoteSafeWithOptions` |

<a id="build-test-and-release-workflow"></a>

## 📦 Build, test, and release workflow / 构建、测试与发布工作流

Clone the source code.<br>
克隆源码：

```bash
git clone https://github.com/imajinyun/go-knifer.git
cd go-knifer
```

Run tests.<br>
运行测试：

```bash
make test
```

Run the CI test-job gates locally. This verifies modules, vet, tidy/diff cleanliness, architecture rules, race/shuffle tests, coverage gates, and the exported API snapshot.

本地运行 CI 测试任务门禁。该命令会校验模块、vet、tidy/diff 清洁度、架构规则、race/shuffle 测试、覆盖率门禁与导出 API 快照。

```bash
make ci-test
```

Run the same local safety checks used by CI before opening a PR.<br>
提交 PR 前运行与 CI 对齐的本地安全检查：

```bash
make check
```

`make check` includes the `ci-test` class of checks plus `golangci-lint` and `govulncheck`.<br>
`make check` 包含 `ci-test` 类检查，并额外运行 `golangci-lint` 和 `govulncheck`。

Benchmark baselines / 基准测试基线：

```bash
make bench-core
make bench-facade
make bench-core BENCHCOUNT=10 BENCHTIME=3s
```

Refresh the API snapshot after an intentional exported API change.<br>
有意修改导出 API 后刷新 API 快照：

```bash
UPDATE_API=1 make api-check
```

<a id="governance"></a>

## 🛡️ Governance / 治理

- Security reports / 安全报告：see [`../../SECURITY.md`](../../SECURITY.md). Please do not disclose suspected vulnerabilities in public issues. / 请查看 [`../../SECURITY.md`](../../SECURITY.md)，不要在公开 issue 中披露疑似漏洞。
- Release notes / 发布说明：see [`../../CHANGELOG.md`](../../CHANGELOG.md). User-visible changes should be recorded before tagging a release. / 用户可见变更应在发布 tag 前记录到 [`../../CHANGELOG.md`](../../CHANGELOG.md)。
- Coverage gate / 覆盖率门禁：CI enforces the repository baseline with `bash bin/check_coverage.sh coverage.out`. Raise thresholds only after adding supporting tests. / CI 通过 `bash bin/check_coverage.sh coverage.out` 强制仓库基线；只有在补充测试后才提高阈值。
- API gate / API 门禁：`make api-check` compares public signatures against [`../api/exports.txt`](../api/exports.txt). / `make api-check` 会将公开签名与 [`../api/exports.txt`](../api/exports.txt) 比较。
- Workflow gates / 工作流门禁：use `make quick-check`, `make security-check`, `make full-check COVERAGE_FILE=/tmp/go-knifer-coverage.out`, and `make ci-test`. / 使用 `make quick-check`、`make security-check`、`make full-check COVERAGE_FILE=/tmp/go-knifer-coverage.out` 与 `make ci-test`。
- Benchmark baseline / 基准测试基线：use `make bench-core` or `make bench-facade`; treat output as a baseline unless `benchstat` proves a performance change. / 使用 `make bench-core` 或 `make bench-facade`；除非 `benchstat` 证明性能变化，否则结果仅作为基线。

<a id="contributing"></a>

## 🤝 Contributing / 贡献

If you find a bug or want to request a new utility, please open a GitHub Issue. It is recommended to include:

如果你发现 bug 或希望新增工具，请提交 GitHub Issue，并建议包含：

- Go version and operating system / Go 版本与操作系统；
- `go-knifer` version or commit / `go-knifer` 版本或 commit；
- Minimal reproducible code / 最小可复现代码；
- Expected behavior and actual behavior / 预期行为与实际行为；
- Related error logs or test output / 相关错误日志或测试输出。

Pull requests are welcome. To keep the toolkit stable, please follow these principles where possible:

欢迎提交 Pull Request。为了保持工具集稳定，请尽量遵循以下原则：

1. Add new capabilities to the appropriate `internal/*` implementation package first, then expose public APIs from the corresponding `v*` package; / 先在合适的 `internal/*` 实现包中添加能力，再通过对应 `v*` 包暴露公开 API；
2. Add necessary comments for new or modified public APIs; / 为新增或修改的公开 API 添加必要注释；
3. Add unit tests for core logic and run `go test ./...` before submitting; / 为核心逻辑补充单元测试，并在提交前运行 `go test ./...`；
4. Keep code formatted with `gofmt`; / 使用 `gofmt` 保持格式化；
5. Avoid unnecessary third-party dependencies and prefer the standard library when possible. / 避免不必要的三方依赖，能用标准库时优先使用标准库。

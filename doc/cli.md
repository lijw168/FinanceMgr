# FinanceMgr CLI 使用文档

本文档基于 `src/analysis-server/cli/command` 中的 Cobra 命令定义整理，适合在本地或服务器环境中直接使用。

## 1. CLI 基本信息

- 可执行程序：`ccs`
- 入口文件：`src/analysis-server/cli/main.go`
- 全局参数：
  - `--server-url`：服务地址，例如 `http://127.0.0.1:8080`
  - `--token`：登录后的 access token
  - `--timeout`：网络超时，单位毫秒
  - `--verbose` / `-v`：输出详细请求与响应信息
  - `--admin` / `-a`：启用管理员模式
  - `--help` / `-h`：查看帮助

## 2. 启动方式

```bash
./ccs --server-url http://127.0.0.1:8080 --help
```

如果未传 `--server-url`，程序会读取环境变量：

```bash
export CC_SERVER_URL=http://127.0.0.1:8080
./ccs --help
```

## 3. 全局帮助示例

```bash
./ccs --help
```

输出中会展示所有可用子命令，如：

- `login`
- `logout`
- `status-checkout`
- `company-create`
- `company-list`
- `company-show`
- `company-update`
- `companyGroup-create`
- `companyGroup-list`
- `accSub-create`
- `accSub-list`
- `operator-create`
- `operator-list`
- `voucher-create`
- `voucher-show`
- `vouInfo-list`
- `yearBal-create`
- `yearBal-list`

## 4. 登录与认证命令

### 4.1 login

```bash
./ccs login <name> <password> <companyID>
```

示例：

```bash
./ccs --server-url http://127.0.0.1:8080 login alice 123456 1
```

说明：

- `name`：用户名
- `password`：登录密码
- `companyID`：公司编号
- 登录成功后，返回包含 `accessToken` 的数据，建议保存到 `--token`

### 4.2 logout

```bash
./ccs logout <operatorID>
```

示例：

```bash
./ccs --server-url http://127.0.0.1:8080 --token "<token>" logout 12
```

### 4.3 loginInfo-list

```bash
./ccs loginInfo-list <operatorId>
```

### 4.4 loginInfo-show

```bash
./ccs loginInfo-show <operatorId>
```

### 4.5 status-checkout

```bash
./ccs status-checkout <operatorID>
```

示例：

```bash
./ccs --server-url http://127.0.0.1:8080 --token "<token>" status-checkout 12
```

## 5. 公司管理命令

### 5.1 company-create

```bash
./ccs company-create <companyName> <abbrevName> <corporator> <phone> <e_mail> <companyAdd> <startAccountPeriod>
```

示例：

```bash
./ccs --server-url http://127.0.0.1:8080 --token "<token>" \
  company-create "Demo Company" "DC" "Alice" "13800000000" "alice@example.com" "Shenzhen" 202401
```

### 5.2 company-list

```bash
./ccs company-list
```

可加选项：

```bash
./ccs company-list --column CompanyID --column CompanyName
```

### 5.3 company-show

```bash
./ccs company-show <companyId>
```

### 5.4 company-update

```bash
./ccs company-update <companyId> [--comName "..."] [--abbrevName "..."] [--comAddr "..."] [--corporator "..."] [--email "..."] [--bc "..."] [--accYear 2025]
```

示例：

```bash
./ccs company-update 1 --comName "Updated Company" --accYear 2025
```

### 5.5 companyGroup-associated

```bash
./ccs companyGroup-associated <companyGroupId> <companyId>
```

参数：

- `--isAttach`：是否附加/关联，默认 `true`

### 5.6 companyAccountYearInfo-list

```bash
./ccs companyAccountYearInfo-list <operatorId>
```

## 6. 公司分组命令

### 6.1 companyGroup-create

```bash
./ccs companyGroup-create <groupName> <groupStatus>
```

示例：

```bash
./ccs companyGroup-create "Finance Group" 1
```

### 6.2 companyGroup-list

```bash
./ccs companyGroup-list
```

### 6.3 companyGroup-show

```bash
./ccs companyGroup-show <companyGroupId>
```

### 6.4 companyGroup-update

```bash
./ccs companyGroup-update <companyGroupId> <groupStatus>
```

可附加参数：

```bash
./ccs companyGroup-update 1 1 --comGroupName "New Group Name"
```

### 6.5 companyGroup-delete

```bash
./ccs companyGroup-delete <companyGroupId>
```

## 7. 会计科目命令

### 7.1 accSub-create

```bash
./ccs accSub-create <commonId> <subjectName> <subjectLevel> <companyId> <subjectDirection> <subjectType> <subjectStyle>
```

示例：

```bash
./ccs accSub-create "1001" "现金" 1 1 1 1 "A"
```

说明：

- `commonId`：科目编码前缀/通用编号
- `subjectLevel`：科目层级
- `subjectDirection`：科目方向（借/贷）
- `subjectType`：科目类型
- `subjectStyle`：科目样式

### 7.2 accSub-list

```bash
./ccs accSub-list <companyId>
```

### 7.3 accSub-show

```bash
./ccs accSub-show <subjectId>
```

### 7.4 accSub-update

```bash
./ccs accSub-update <subjectID> <commonID> <subjectName> <subjectLevel> <companyId> <subjectDirection> <subjectType>
```

### 7.5 accSub-createTemplate

```bash
./ccs accSub-createTemplate <companyId>
```

### 7.6 accSub-delete

```bash
./ccs accSub-delete <subjectId>
```

## 8. 操作员命令

### 8.1 operator-create

```bash
./ccs operator-create <companyID> <name> <password> <job> <department> <role>
```

示例：

```bash
./ccs operator-create 1 bob 123456 "财务经理" "财务部" 1
```

### 8.2 operator-list

```bash
./ccs operator-list <companyId>
```

### 8.3 operator-show

```bash
./ccs operator-show <operatorID>
```

### 8.4 operator-update

```bash
./ccs operator-update <operatorID> <job>
```

### 8.5 operator-delete

```bash
./ccs operator-delete <operatorID>
```

## 9. 凭证命令

### 9.1 voucher-create

```bash
./ccs voucher-create <companyID> <voucherMonth> <voucherFiller>
```

可选参数：

- `--subject`：科目名
- `--summary`：摘要
- `--dm`：借方金额
- `--cm`：贷方金额
- `--sub1`：科目1编号
- `--sub2`：科目2编号

示例：

```bash
./ccs voucher-create 1 8 alice --subject "现金" --summary "测试" --dm 100 --cm 0 --sub1 100 --sub2 101
```

### 9.2 voucher-show

```bash
./ccs voucher-show <voucherId> <voucherYear>
```

### 9.3 voucher-delete

```bash
./ccs voucher-delete <voucherId> <voucherYear>
```

### 9.4 voucher-arrange

```bash
./ccs voucher-arrange <companyID> <voucherYear> <voucherMonth>
```

可选参数：

```bash
./ccs voucher-arrange 1 2025 8 --isArrangeVoucherNum true
```

### 9.5 vouRecord-list

```bash
./ccs vouRecord-list <voucherId> <voucherYear>
```

### 9.6 vouInfo-show

```bash
./ccs vouInfo-show <voucherID> <voucherYear>
```

### 9.7 vouInfo-list

```bash
./ccs vouInfo-list <companyId> <voucherYear>
```

### 9.8 vouInfo-getLatest

```bash
./ccs vouInfo-getLatest <companyID> <voucherYear> <voucherMonth>
```

### 9.9 vouInfo-getMaxNumOfMonth

```bash
./ccs vouInfo-getMaxNumOfMonth <companyID> <voucherYear> <voucherMonth>
```

### 9.10 vouInfo-update

```bash
./ccs vouInfo-update <voucherId> <voucherYear>
```

可选参数：

- `--vouFiller`
- `--vouAuditor`
- `--vouDate`
- `--billCount`
- `--status`

## 10. 年度余额命令

### 10.1 yearBal-create

```bash
./ccs yearBal-create <companyID> <year> <subjectID>
```

可选参数：

```bash
./ccs yearBal-create 1 2025 100 --balance 1234.56
```

### 10.2 yearBal-show

```bash
./ccs yearBal-show <companyId> <year> <subjectId>
```

### 10.3 yearBal-update

```bash
./ccs yearBal-update <companyID> <year> <subjectID>
```

可选参数：

```bash
./ccs yearBal-update 1 2025 100 --balance 4567.89 --status 1
```

### 10.4 yearBal-list

```bash
./ccs yearBal-list <companyId> <year> <subjectId>
```

### 10.5 yearBal-delete

```bash
./ccs yearBal-delete <companyId> <year> <subjectId>
```

### 10.6 yearBal-accSubBal-show

```bash
./ccs yearBal-accSubBal-show <companyId> <year> <subjectId>
```

## 11. 常用最小工作流

### 11.1 典型登录与查询流程

```bash
./ccs --server-url http://127.0.0.1:8080 login alice 123456 1
./ccs --server-url http://127.0.0.1:8080 --token "<accessToken>" company-list
./ccs --server-url http://127.0.0.1:8080 --token "<accessToken>" accSub-list 1
```

### 11.2 创建公司与科目

```bash
./ccs --server-url http://127.0.0.1:8080 --token "<token>" \
  company-create "Demo Company" "DC" "Alice" "13800000000" "alice@example.com" "Shenzhen" 202401

./ccs --server-url http://127.0.0.1:8080 --token "<token>" \
  accSub-create "1001" "现金" 1 1 1 1 "A"
```

### 11.3 创建凭证

```bash
./ccs --server-url http://127.0.0.1:8080 --token "<token>" \
  voucher-create 1 8 alice --subject "现金" --summary "测试" --dm 100 --cm 0 --sub1 100 --sub2 101
```

## 12. 注意事项

1. 很多命令的参数是基于数字 ID 的，调用前应先查询对应对象的列表。
2. `--token` 需要在登录成功后使用，未登录时很多接口会返回“please login first.”。
3. `--verbose` 可用于排查请求参数和响应错误，适合故障定位。
4. `--admin` 仅适用于需要管理员权限的场景。
5. CLI 实际是通过 SDK 向 `?Action=` 形式的 HTTP 接口发起 POST 请求，因此与直接调用后台接口协议保持一致。

## 13. 代码定位

本文档主要依据以下文件：

- `src/analysis-server/cli/main.go`
- `src/analysis-server/cli/command/init.go`
- `src/analysis-server/cli/command/*.go`
- `src/analysis-server/sdk/util/request.go`
- `src/analysis-server/sdk/mgr/*.go`

如果需要进一步扩展文档，可以继续补充每个命令的输出字段说明和 JSON 示例。

# FinanceMgr API 接口文档

本文档基于代码中的 HTTP 路由注册、SDK 调用和参数模型整理，覆盖 `analysis-server` 的主要接口能力。

## 1. 通用协议

- 服务入口：`http://<host>:<port>`
- 通信方式：HTTP POST
- 调用方式：在 URL 中追加 `?Action=<ActionName>`
- 请求头：
  - `Content-Type: application/json`
  - 认证：`Cookie: access_token=<token>`
  - 管理员模式：`Secret-Id: MgrClientSecretId`
- 返回格式：统一 JSON

通用返回结构：

```json
{
  "code": 0,
  "message": "success",
  "detail": "",
  "data": {}
}
```

其中：

- `code`：0 表示成功，非 0 表示错误
- `message`：简短错误/成功信息
- `detail`：更多细节信息
- `data`：实际返回对象或列表

通用分页返回：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_count": 123,
    "elements": [
      { "companyId": 1, "companyName": "Demo" }
    ]
  }
}
```

## 2. 接口入口与路由规则

所有业务接口都由 `src/common/url/urlrouter.go` 统一分发，实际请求是：

```http
POST /?Action=CreateCompany
Content-Type: application/json
```

body 中是对应的 JSON 参数对象。

还有一类公共接口由 `src/common/url/default_urlrouter.go` 注册：

- `setloglevel`：设置日志级别，参数 `level`（trace/debug/info/warn/error/fatal）
- `getloglevel`：读取日志级别
- `version`：获取版本信息

## 3. 参数模型说明

本系统的参数模型主要在 `src/analysis-server/model/params.go` 中定义，最常见的共有几类：

- `BaseParams`：`{"id": 123}`
- `ListParams`：分页和筛选
- `FilterItem`：筛选条件
- `OrderItem`：排序条件
- 业务参数对象：例如 `CreateCompanyParams`、`CreateVoucherParams`、`AuthenInfoParams`

### 3.1 常见筛选参数

```json
{
  "filter": [
    { "field": "companyId", "value": 1 }
  ],
  "orders": [
    { "field": "createdAt", "direction": 1 }
  ],
  "desc_offset": 0,
  "desc_limit": 20
}
```

说明：

- `field`：字段名
- `value`：筛选值
- `direction`：排序方向，通常 1 表示升序，-1 表示降序
- `desc_offset` / `desc_limit`：分页偏移和数量

## 4. API 行为总览

以下按模块列出主要 `Action` 名称与职责。

### 4.1 公司分组（CompanyGroup）

- `CreateCompanyGroup`：创建公司分组
- `DeleteCompanyGroup`：删除公司分组
- `GetCompanyGroup`：获取分组详情
- `ListCompanyGroup`：分页/条件查询分组列表
- `UpdateCompanyGroup`：更新分组信息

### 4.2 公司（Company）

- `CreateCompany`：创建公司
- `DeleteCompany`：删除公司
- `GetCompany`：查询公司详情
- `ListCompany`：查询公司列表
- `UpdateCompany`：更新公司信息
- `AssociatedCompanyGroup`：关联/取消关联公司分组
- `ListCompanyAccountYearInfo`：列出公司建账年信息

### 4.3 会计科目（Account Subject）

- `CreateAccSub`：创建会计科目
- `DeleteAccSub`：删除会计科目
- `ListAccSub`：查询科目列表
- `GetAccSub`：查询科目详情
- `UpdateAccSub`：更新科目
- `QueryAccSubReference`：查询科目引用关系
- `CopyAccSubTemplate`：复制科目模板
- `GenerateAccSubTemplate`：生成科目模板

### 4.4 年度余额（Year Balance）

- `GetYearBalance`：查询某年度某科目的余额
- `GetAccSubYearBalValue`：查询单科目年余额值
- `CreateYearBalance`：创建年余额记录
- `BatchCreateYearBalance`：批量创建年余额记录
- `UpdateYearBalance`：更新年余额记录
- `BatchUpdateBals`：批量更新余额
- `DeleteYearBalance`：删除年余额记录
- `BatchDeleteYearBalance`：批量删除年余额记录
- `ListYearBalance`：查询年余额列表
- `AnnualClosing`：年度结账
- `CancelAnnualClosing`：取消结账
- `GetAnnualClosingStatus`：查询结账状态

### 4.5 操作员（Operator）

- `CreateOperator`：创建操作员
- `DeleteOperator`：删除操作员
- `GetOperatorInfo`：获取操作员信息
- `ListOperatorInfo`：查询操作员列表
- `UpdateOperator`：更新操作员信息

### 4.6 认证与登录（Auth）

- `Login`：登录
- `Logout`：登出
- `StatusCheckout`：检查用户登录状态
- `ListLoginInfo`：查询登录日志/登录信息

### 4.7 凭证（Voucher）

- `CreateVoucher`：新建凭证
- `UpdateVoucher`：更新凭证
- `DeleteVoucher`：删除凭证
- `ArrangeVoucher`：整理凭证编号/顺序
- `ListVoucherRecords`：查询凭证分录列表
- `GetVoucherInfo`：获取凭证主信息
- `GetVoucher`：获取凭证完整信息（主信息 + 明细）
- `GetLatestVoucherInfo`：获取指定月最新凭证信息
- `ListVoucherInfo`：分页查询凭证主信息
- `ListVoucherInfoWithAuxCondition`：按附加条件查询凭证主信息
- `GetMaxNumOfMonth`：查询某月最大凭证号
- `UpdateVoucherInfo`：更新凭证信息
- `BatchAuditVouchers`：批量审核凭证
- `CalculateAccumulativeMoney`：计算累计金额
- `BatchCalcAccuMoney`：批量计算累计金额
- `CalcAccountOfPeriod`：按期间计算科目金额
- `GetNoAuditedVoucherInfoCount`：查询未审核凭证数

### 4.8 模板与菜单

- `CreateVoucherTemplate`：创建凭证模板
- `DeleteVoucherTemplate`：删除凭证模板
- `GetVoucherTemplate`：获取凭证模板
- `ListVoucherTemplate`：查询模板列表
- `ListMenuInfo`：查询菜单信息

## 5. 具体接口说明

### 5.1 登录与认证

#### Login

- Action：`Login`
- 请求体：`AuthenInfoParams`

```json
{
  "name": "alice",
  "password": "123456",
  "companyId": 1
}
```

返回值通常为登录用户信息与 `accessToken`，并写入 cookie。

#### Logout

- Action：`Logout`
- 请求体：`{"id": 12}`

#### StatusCheckout

- Action：`StatusCheckout`
- 请求体：`{"id": 12}`

#### ListLoginInfo

- Action：`ListLoginInfo`
- 请求体：`ListParams` + `filter[operatorId]`

---

### 5.2 公司管理

#### CreateCompany

- Action：`CreateCompany`
- 请求体：`CreateCompanyParams`

```json
{
  "companyName": "Demo Company",
  "abbreviationName": "DC",
  "corporator": "Alice",
  "phone": "13800000000",
  "e_mail": "alice@example.com",
  "companyAddr": "Shenzhen",
  "backup": "",
  "startAccountPeriod": 202401
}
```

#### UpdateCompany

- Action：`UpdateCompany`
- 请求体：`ModifyCompanyParams`

```json
{
  "companyId": 1,
  "companyName": "Updated Company",
  "abbreviationName": "UC",
  "latestAccountYear": 2025
}
```

#### ListCompany

- Action：`ListCompany`
- 请求体：`ListParams`

```json
{
  "filter": [
    { "field": "companyName", "value": "Demo" }
  ],
  "orders": [
    { "field": "updatedAt", "direction": -1 }
  ],
  "desc_offset": 0,
  "desc_limit": 10
}
```

#### GetCompany

- Action：`GetCompany`
- 请求体：`{"id": 1}`

#### AssociatedCompanyGroup

- Action：`AssociatedCompanyGroup`
- 请求体：

```json
{
  "companyGroupId": 1,
  "companyId": 2,
  "isAttach": true
}
```

---

### 5.3 公司分组管理

#### CreateCompanyGroup

- Action：`CreateCompanyGroup`
- 请求体：

```json
{
  "groupName": "Finance Group",
  "groupStatus": 1
}
```

#### UpdateCompanyGroup

- Action：`UpdateCompanyGroup`
- 请求体：

```json
{
  "companyGroupId": 1,
  "groupName": "Updated Group",
  "groupStatus": 1
}
```

---

### 5.4 会计科目管理

#### CreateAccSub

- Action：`CreateAccSub`
- 请求体：`CreateSubjectParams`

```json
{
  "companyId": 1,
  "commonId": "1001",
  "subjectName": "现金",
  "subjectLevel": 1,
  "subjectDirection": 1,
  "subjectType": 1,
  "subjectStyle": "A"
}
```

说明：

- `subjectDirection`：借/贷方向
- `subjectType`：科目类别
- `subjectStyle`：科目样式标识

#### UpdateAccSub

- Action：`UpdateAccSub`
- 请求体：`ModifySubjectParams`

#### ListAccSub

- Action：`ListAccSub`
- 请求体：`ListParams`

```json
{
  "filter": [
    { "field": "companyId", "value": 1 }
  ],
  "desc_offset": 0,
  "desc_limit": 50
}
```

#### GetAccSub

- Action：`GetAccSub`
- 请求体：`{"id": 123}`

---

### 5.5 年度余额管理

#### CreateYearBalance

- Action：`CreateYearBalance`
- 请求体：

```json
{
  "companyId": 1,
  "year": 2025,
  "subjectId": 100,
  "balance": 1234.56,
  "status": 0
}
```

#### BatchCreateYearBalance

- Action：`BatchCreateYearBalance`
- 请求体：

```json
{
  "companyId": 1,
  "year": 2025,
  "optSubAndBals": [
    { "subjectId": 100, "balance": 1000 },
    { "subjectId": 101, "balance": 2000 }
  ]
}
```

#### ListYearBalance

- Action：`ListYearBalance`
- 请求体：`ListParams` + `filter[companyId]`

#### AnnualClosing

- Action：`AnnualClosing`
- 请求体：`{"companyId": 1, "year": 2025}`

---

### 5.6 操作员管理

#### CreateOperator

- Action：`CreateOperator`
- 请求体：

```json
{
  "companyId": 1,
  "name": "bob",
  "password": "123456",
  "job": "财务经理",
  "department": "财务部",
  "Status": 1,
  "role": 1
}
```

#### ListOperatorInfo

- Action：`ListOperatorInfo`
- 请求体：`ListParams`

```json
{
  "filter": [
    { "field": "companyId", "value": 1 }
  ],
  "desc_offset": 0,
  "desc_limit": 20
}
```

---

### 5.7 凭证管理

#### CreateVoucher

- Action：`CreateVoucher`
- 请求体：

```json
{
  "infoParams": {
    "companyId": 1,
    "voucherDate": 15,
    "voucherFiller": "alice",
    "billCount": 2
  },
  "recordsParams": [
    {
      "voucherId": 0,
      "subjectName": "现金",
      "debitMoney": 100.00,
      "creditMoney": 0,
      "summary": "测试",
      "subId1": 100
    },
    {
      "voucherId": 0,
      "subjectName": "应收账款",
      "debitMoney": 0,
      "creditMoney": 100.00,
      "summary": "测试",
      "subId1": 101
    }
  ]
}
```

#### GetVoucher

- Action：`GetVoucher`
- 请求体：

```json
{
  "voucherYear": 2025,
  "id": 42
}
```

#### ListVoucherRecords

- Action：`ListVoucherRecords`
- 请求体：`ListParams` + `filter[voucherId]` + `filter[voucherYear]`

#### BatchAuditVouchers

- Action：`BatchAuditVouchers`
- 请求体：

```json
{
  "voucherYear": 2025,
  "ids": [1, 2, 3],
  "status": 1,
  "voucherAuditor": "alice"
}
```

#### ArrangeVoucher

- Action：`ArrangeVoucher`
- 请求体：

```json
{
  "voucherYear": 2025,
  "companyId": 1,
  "voucherMonth": 8,
  "isArrangeVoucherNum": true
}
```

---

## 6. 常见请求示例

### 6.1 cURL 示例：登录

```bash
curl -X POST "http://127.0.0.1:8080/?Action=Login" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "alice",
    "password": "123456",
    "companyId": 1
  }'
```

### 6.2 cURL 示例：查询公司列表

```bash
curl -X POST "http://127.0.0.1:8080/?Action=ListCompany" \
  -H "Content-Type: application/json" \
  -d '{
    "filter": [],
    "orders": [],
    "desc_offset": 0,
    "desc_limit": 10
  }'
```

### 6.3 cURL 示例：创建凭证

```bash
curl -X POST "http://127.0.0.1:8080/?Action=CreateVoucher" \
  -H "Content-Type: application/json" \
  -d '{
    "infoParams": {
      "companyId": 1,
      "voucherDate": 15,
      "voucherFiller": "alice",
      "billCount": 2
    },
    "recordsParams": [
      {
        "voucherId": 0,
        "subjectName": "现金",
        "debitMoney": 100,
        "creditMoney": 0,
        "summary": "测试",
        "subId1": 100
      }
    ]
  }'
```

## 7. 代码定位

本文档主要基于以下代码：

- `src/analysis-server/api/main/registerHandler.go`
- `src/analysis-server/model/params.go`
- `src/analysis-server/model/view.go`
- `src/common/url/urlrouter.go`
- `src/analysis-server/sdk/util/request.go`

## 8. 实施建议

1. 生产环境下，建议统一 API 网关入口，保持 `Action` 参数方式不变。
2. 请求体如果较复杂，优先配合 SDK 中 `options` 结构体和 `model` 参数结构来生成请求。
3. 若要自动化文档生成，可从 `registerHandler.go` 中读取 `Action` 名称，并将其映射到对应 `params` / `view` 结构体。

以上内容可作为 FinanceMgr 的基础接口文档，后续可以继续按模块扩展为 OpenAPI/Swagger 版本。

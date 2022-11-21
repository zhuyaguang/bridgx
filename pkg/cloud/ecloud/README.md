# ECloud 移动云

- [OpenAPI Explorer](https://ecloud.10086.cn/op-oneapi-static/#/overview)
- [SDK 中心](https://ecloud.10086.cn/op-oneapi-static/#/center/sdk)
- [Go SDK 使用指南](https://ecloud.10086.cn/op-help-center/doc/article/53799)


- [对象存储 EOS API文档](https://ecloud.10086.cn/op-help-center/doc/article/40960)
- [对象存储 SDK](https://ecloud.10086.cn/op-help-center/doc/article/24569)

### 注意
- 移动云所有接口调用应先判断 `err`，再判断 `*Model.xxxResponse`，否则可能会忽略掉一些业务异常。因为 `err` 一般是调用链路上的错误，如果调用成功一般 `err == nil`； `Response` 代表请求正确处理，但是业务处理可能有异常。
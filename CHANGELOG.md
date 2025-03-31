## [](https://github.com/kysion/base-library/compare/v0.3.0...v) (2025-03-31)

### Bug Fixes

* **utility/kmap:** 修正 HashMap 类型的 MarshalJSON 方法接收者 ([eecc67d](https://github.com/kysion/base-library/commit/eecc67d59fc0a79afd8c0a3758d8429193890a90))
## [0.3.0](https://github.com/kysion/base-library/compare/v0.2.4...v0.3.0) (2025-03-29)

### Features

* **base_tree:** 为树结构添加子节点排序功能 ([9cf57f5](https://github.com/kysion/base-library/commit/9cf57f5878e7b82816ab083a73564f4549251091))
## [0.2.4](https://github.com/kysion/base-library/compare/v0.2.3...v0.2.4) (2025-03-23)

### Features

* **enum:** 优化枚举工具包并添加性能测试结果 ([dbde901](https://github.com/kysion/base-library/commit/dbde90111608a4ea10da431d5acd8c16d8761ce8))
* **utility:** 添加枚举工具包 ([564e36f](https://github.com/kysion/base-library/commit/564e36fab3028a45a0c070aff054a00c537ff649))
## [0.2.0](https://github.com/kysion/base-library/compare/v0.1.9...v0.2.0) (2025-01-15)

### Features

* .env 追加多环境配置自定义及隔离支持 ([5db3294](https://github.com/kysion/base-library/commit/5db3294412c92536d5c6309b372635f702644302))
* 增加 MakeModel 方法，支持将查询参数转成 Model ([7a08597](https://github.com/kysion/base-library/commit/7a085972761cf596f617f650697d953158e9bd18))
* 增加 MakeModel 方法，支持将查询参数转成 Model ([23ed361](https://github.com/kysion/base-library/commit/23ed361ac40bc920502a7542b42dcb0d773fccbe))
* 增加工具方法 ([d29d037](https://github.com/kysion/base-library/commit/d29d037ac01bd570f766a4f9a89b00aa5e8601ee))
* 新增切片操作函数 ([18a39e4](https://github.com/kysion/base-library/commit/18a39e453b934efdceffdc8fca32e33830b8fe4d))
* 添加dao配套模板示例文件 ([c1ee2e8](https://github.com/kysion/base-library/commit/c1ee2e8c507b4b549dd44d8f48519a7336b30926))
* 添加工具方法 ([b8985bd](https://github.com/kysion/base-library/commit/b8985bd01dbff50797a73c82dec0d759cea59da7))
* 添加工具方法 ([7956f81](https://github.com/kysion/base-library/commit/7956f811e26f199424dd27e5ba3cd75de567e9a0))
* 添加根据身份证分析关键数据的方法 ([c66d02b](https://github.com/kysion/base-library/commit/c66d02bf80a8c526c3cf4c52f50a82399e691626))
* 添加跨进程通信的解决方案代码（待完善） ([d3f841c](https://github.com/kysion/base-library/commit/d3f841ce771a40f2ec4b07e343185e86042fabe5))
* 补充完整跨进程解决方案 ([953f740](https://github.com/kysion/base-library/commit/953f7407255164af4f3dfe87227dd19f9ed71d21))

### Bug Fixes

* 修复枚举定义的错误用法 ([c0431c1](https://github.com/kysion/base-library/commit/c0431c16f9665eefab703996c2212422dd9857cf))
* 修复查询分页Bug ([6a70da5](https://github.com/kysion/base-library/commit/6a70da53f93e4a82f7d20b75eed3357eafb476d4))
* 修复缓存命中条件逻辑BUG ([788b3ad](https://github.com/kysion/base-library/commit/788b3ad6cbdfdb66b00be68789e64bff64d26069))
* 修复邮箱IsEmail 方法不能校验 点 “ . ”  的Bug ([a6d7e2d](https://github.com/kysion/base-library/commit/a6d7e2df8d13f01413360d57d936d70febc3bdd0))
* 增强应用便捷性 ([4baa23b](https://github.com/kysion/base-library/commit/4baa23b0b6c224aaab2b91b28aa92c90c4e23661))

### Performance Improvements

* 优化命名规范，优化接口注释 ([3b3e901](https://github.com/kysion/base-library/commit/3b3e90177621ac350f9bef6545633c901e42b094))
* 优化命名规范，优化注释说明 ([85240ba](https://github.com/kysion/base-library/commit/85240baea5b25e9014d59686440749eef4eab65b))
* 优化扩展ORM行为条件应用和实现 ([07077b1](https://github.com/kysion/base-library/commit/07077b15591d7ca2081bb468c2c87c05d04f7c23))
* 优化描述信息，优化代码格式 ([9d580c3](https://github.com/kysion/base-library/commit/9d580c3112b9468ca95892471cdfde530d4acdc8))
* 优化逻辑，增加可读性 ([7a47aa5](https://github.com/kysion/base-library/commit/7a47aa52b2b5621dd34aded81d114d30a7b1c2bc))
* 增强ORM缓存的控制，支持仅忽略同一个上下文相关表的缓存 ([a4bf28e](https://github.com/kysion/base-library/commit/a4bf28e7f9d3d2c5359aee2e332ad44605a35aaf))
* 提升效率、增加线程安全、增加可读性、优化大文件防止内存占用高的问题 ([692c6d0](https://github.com/kysion/base-library/commit/692c6d0854df77ed4e5897e03925adadb07ce662))
* 新增数据缓存忽略控制支持、新增model 前置条件附加支持、丰富逻辑实现相关注释 ([e5ba9b2](https://github.com/kysion/base-library/commit/e5ba9b2f2c6fe28de1cc6e17e1dc70357b7383e4))
* 移除对idgen的依赖，增加工厂方法来初始化权限工厂方法对象，增加工厂初始化校验逻辑，增加代码可读性 ([7876b17](https://github.com/kysion/base-library/commit/7876b172f5fde2bb6cdf3d1a1563cdc5a45cc883))
## [0.0.33](https://github.com/kysion/base-library/compare/v0.0.32...v0.0.33) (2024-05-11)

### Features

* 添加部分方法封装 ([20ff438](https://github.com/kysion/base-library/commit/20ff43877b02540038ab8646158f1bce742ac89e))
## [0.0.32](https://github.com/kysion/base-library/compare/v0.0.31...v0.0.32) (2024-02-23)

### Features

* 添加防抖函数 ([014a4f4](https://github.com/kysion/base-library/commit/014a4f4ec96b21c864d8512c71d2ffa46fe772a0))

### Bug Fixes

* 修复SQL查询通用模块中Count的Bug ([068812d](https://github.com/kysion/base-library/commit/068812dbf540fd513cecaac270c96a8521d75489))
* 修复过滤拓展字段方法的逻辑 ([5ca823d](https://github.com/kysion/base-library/commit/5ca823d18e929f813b51309c86430694af7bdb4b))
## [0.0.30](https://github.com/kysion/base-library/compare/v0.0.29...v0.0.30) (2023-12-21)

### Features

* 添加配置示例文件，使用最新规范 ([67dd3ca](https://github.com/kysion/base-library/commit/67dd3ca048c399d647656d4550ce8fc6756d02d7))

### Bug Fixes

* 修复swaggerPath错误Bug ([358f0e2](https://github.com/kysion/base-library/commit/358f0e2ceb4d819d0cbb9d73fb3112fc037f0b13))
## [0.0.28](https://github.com/kysion/base-library/compare/v0.0.27...v0.0.28) (2023-11-20)

### Features

* 添加封装方法, 判断指定时间是一周内的第几天 ([7f5d3cf](https://github.com/kysion/base-library/commit/7f5d3cf46f389fc03177e714de4b429a963ac488))
* 添加工具类fsql ([ba329f4](https://github.com/kysion/base-library/commit/ba329f43a05980eac2641d194a4f9d0f8fa3dc18))

### Performance Improvements

* 优化权限工程定义 ([1661d97](https://github.com/kysion/base-library/commit/1661d97a923ba0314de889159e52e973dc5b9a10))
* 枚举新增 Has 方法，判断是否包含指定枚举值 ([5f516e5](https://github.com/kysion/base-library/commit/5f516e54c7ea1f318f471e5e026d00f54127c106))
* 枚举类型新增 Add、Remove 方法 ([9a88356](https://github.com/kysion/base-library/commit/9a883568aac2c504eddf0a6c4cdf73a78a3aaf3a))
## [0.0.21](https://github.com/kysion/base-library/compare/v0.0.20...v0.0.21) (2023-07-29)

### Features

* 将验证码 枚举抽取到base库 ([b0f6564](https://github.com/kysion/base-library/commit/b0f6564fb0ebf6050578ce9ea3b276c1047783f5))
* 新增常用类型校验判断方法 ([3ce2488](https://github.com/kysion/base-library/commit/3ce2488de48c4f683cdc1896f8a08b1130e3b5ba))
* 添加SecondToDuration 方法 ([17ca4cd](https://github.com/kysion/base-library/commit/17ca4cd17e04d2bc0f6e31ab320cff83ddc29651))
* 添加文件获取字节的方法 ([a704e41](https://github.com/kysion/base-library/commit/a704e41063e51f9cb42fc04fc13019a425131552))
* 添加时间季度处理方法 ([cf08c00](https://github.com/kysion/base-library/commit/cf08c00369a86f11533ec549b0389e7539e68a21))
* 添加是否是手机号和邮箱方法 ([bdcc196](https://github.com/kysion/base-library/commit/bdcc196a474ec0fc5f67b6e1693c004d080c4b3d))
* 重构树构建逻辑，实现毫秒级完成列表转树的构建 ([7e1682f](https://github.com/kysion/base-library/commit/7e1682fe4a739574024d233cb148141fd1dd1aac))
* 验证码枚举添加设置邮箱 ([db1c8b2](https://github.com/kysion/base-library/commit/db1c8b28fcb2c7e15ed2c50373ca6429cf24160c))
* 验证码类型添加忘记用户名&密码类型, 此类型的验证码支持复用 ([c4d2155](https://github.com/kysion/base-library/commit/c4d2155daa7ad5c1b9469a27937f126625e671bb))

### Bug Fixes

* 修复base-hook被引用空指针错误 ([e9473fa](https://github.com/kysion/base-library/commit/e9473fa1ae8ca7239a5127a732108029549e8d8f))
* 修复项目在Goland 中结构错误的问题 ([59db12b](https://github.com/kysion/base-library/commit/59db12b0b61eabffc36c8120f7c40198be97531f))

### Performance Improvements

* 优化枚举类型应用规范 ([7bec7eb](https://github.com/kysion/base-library/commit/7bec7eb287046498cd5093443de03eecc900f757))
## [0.0.3](https://github.com/kysion/base-library/compare/v0.0.1...v0.0.3) (2023-02-19)

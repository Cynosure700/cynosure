启动新的代理来自主处理复杂的多步骤任务。

可用的代理类型及其可访问的工具：

- general-purpose：用于研究复杂问题、搜索代码和执行多步骤任务的通用代理。当您搜索关键词或文件并且不确定在前几次尝试中能找到正确匹配时，请使用此代理为您执行搜索。（工具：\_）
- statusline-setup：使用此代理配置用户的Claude Code状态栏设置。（工具：Read, Edit）
- output-style-setup：使用此代理创建Claude Code输出样式。（工具：Read, Write, Edit, Glob, LS, Grep）
- sales-automator：起草冷邮件、后续邮件和提案模板。创建定价页面、案例研究和销售脚本。主动用于销售外展或客户开发。（工具：\_）
- architect-reviewer：审查代码更改的架构一致性和模式。在任何结构性更改、新服务或API修改后主动使用。确俚SOLID原则、适当的分层和可维护性。（工具：\_）
- customer-support：处理支持工单、FAQ回复和客户邮件。创建帮助文档、故障排除指南和模板回复。主动用于客户查询或支持文档。（工具：\_）
- ai-engineer：构建LLM应用、RAG系统和提示管道。实现向量搜索、代理编排和AI API集成。主动用于LLM功能、聊天机器人或AI驱动的应用程序。（工具：\_）
- database-optimizer：优化SQL查询、设计高效索引和处理数据库迁移。解决N+1问题、慢查询和实现缓存。主动用于数据库性能问题或模式优化。（工具：\_）
- legal-advisor：起草隐私政策、服务条款、免责声明和法律通知。创建符合GDPR的文本、cookie政策和数据处理协议。主动用于法律文档、合规文本或监管要求。（工具：\_）
- search-specialist：使用高级搜索技术和综合的专家网络研究员。精通搜索操作符、结果过滤和多源验证。处理竞争分析和事实核查。主动用于深度研究、信息收集或趋势分析。（工具：\_）
- tutorial-engineer：从代码创建分步教程和教育内容。将复杂概念转化为渐进式学习体验，带有实践示例。主动用于入职指南、功能教程或概念解释。（工具：\_）
- error-detective：在日志和代码库中搜索错误模式、堆栈跟踪和异常。关联跨系统错误并识别根本原因。主动用于调试问题、分析日志或调查生产错误。（工具：\_）
- mlops-engineer：构建ML管道、实验跟踪和模型注册表。实现MLflow、Kubeflow和自动重训。处理数据版本控制和可重现性。主动用于ML基础设施、实验管理或管道自动化。（工具：\_）
- database-admin：管理数据库操作、备份、复制和监控。处理用户权限、维护任务和灾难恢复。主动用于数据库设置、操作问题或恢复程序。（工具：\_）
- rust-pro：编写具有所有权模式、生命周期和trait实现的惯用Rust代码。精通async/await、安全并发和零成本抽象。主动用于Rust内存安全、性能优化或系统编程。（工具：\_）
- risk-manager：监控投资组合风险、R倍数和位置限制。创建对冲策略、计算期望值和实施止损。主动用于风险评估、交易跟踪或投资组合保护。（工具：\_）
- minecraft-bukkit-pro：精通使用Bukkit、Spigot和Paper API进行Minecraft服务器插件开发。专门从事事件驱动架构、命令系统、世界操作、玩家管理和性能优化。主动用于插件架构、游戏机制、服务器端功能或跨版本兼容性。（工具：\_）
- quant-analyst：构建金融模型、回测交易策略和分析市场数据。实现风险指标、投资组合优化和统计套利。主动用于量化金融、交易算法或风险分析。（工具：\_）
- api-documenter：创建OpenAPI/Swagger规范、生成SDK和编写开发者文档。处理版本控制、示例和交互式文档。主动用于API文档或客户端库生成。（工具：\_）
- cpp-pro：编写具有现代特性、RAII、智能指针和STL算法的惯用C++代码。处理模板、移动语义和性能优化。主动用于C++重构、内存安全或复杂C++模式。（工具：\_）
- performance-engineer：分析应用程序、优化瓶颈和实现缓存策略。处理负载测试、CDN设置和查询优化。主动用于性能问题或优化任务。（工具：\_）
- debugger：专门调试错误、测试失败和意外行为的调试专家。在遇到任何问题时主动使用。（工具：\_）
- legacy-modernizer：重构遗留代码库、迁移过时框架和实现渐进式现代化。处理技术债务、依赖更新和向后兼容性。主动用于遗留系统更新、框架迁移或技术债务减少。（工具：\_）
- golang-pro：编写具有goroutine、channel和interface的惯用Go代码。优化并发、实现Go模式和确保适当的错误处理。主动用于Go重构、并发问题或性能优化。（工具：\_）
- docs-architect：从现有代码库创建全面的技术文档。分析架构、设计模式和实现细节，以生成长篇技术手册和电子书。主动用于系统文档、架构指南或技术深度分析。（工具：\_）
- test-automator：创建包含单元、集成和e2e测试的全面测试套件。设置CI管道、模拟策略和测试数据。主动用于测试覆盖率改进或测试自动化设置。（工具：\_）
- ml-engineer：实现ML管道、模型服务和特征工程。处理TensorFlow/PyTorch部署、A/B测试和监控。主动用于ML模型集成或生产部署。（工具：\_）
- prompt-engineer：为LLM和AI系统优化提示。在构建AI功能、改进代理性能或制作系统提示时使用。提示模式和技术专家。（工具：\_）
- java-pro：精通具有流、并发和JVM优化的现代Java。处理Spring Boot、响应式编程和企业模式。主动用于Java性能调优、并发编程或复杂企业解决方案。（工具：\_）
- scala-pro：精通企业级Scala开发，包括函数式编程、分布式系统和大数据处理。Apache Pekko、Akka、Spark、ZIO/Cats Effect和响应式架构专家。主动用于Scala系统设计、性能优化或企业集成。（工具：\_）
- terraform-specialist：编写高级Terraform模块、管理状态文件和实现IaC最佳实践。处理提供者配置、工作区管理和漂移检测。主动用于Terraform模块、状态问题或IaC自动化。（工具：\_）
- dx-optimizer：开发者体验专家。改进工具、设置和工作流程。在设置新项目、团队反馈后或注意到开发摩擦时主动使用。（工具：\_）
- ios-developer：使用Swift/SwiftUI开发原生iOS应用程序。精通UIKit/SwiftUI、Core Data、网络和应用生命周期。主动用于iOS特定功能、App Store优化或原生iOS开发。（工具：\_）
- code-reviewer：专家代码审查专家。主动审查代码的质量、安全性和可维护性。在编写或修改代码后立即使用。（工具：\_）
- deployment-engineer：配置CI/CD管道、Docker容器和云部署。处理GitHub Actions、Kubernetes和基础设施自动化。在设置部署、容器或CI/CD工作流程时主动使用。（工具：\_）
- backend-architect：设计RESTful API、微服务边界和数据库模式。审查系统架构的可扩展性和性能瓶颈。在创建新的后端服务或API时主动使用。（工具：\_）
- elixir-pro：编写具有OTP模式、监督树和Phoenix LiveView的惯用Elixir代码。精通并发、容错和分布式系统。主动用于Elixir重构、OTP设计或复杂BEAM优化。（工具：\_）
- reference-builder：创建详尽的技术参考和API文档。生成全面的参数列表、配置指南和可搜索的参考资料。主动用于API文档、配置参考或完整的技术规范。（工具：\_）
- devops-troubleshooter：调试生产问题、分析日志和修复部署失败。精通监控工具、事件响应和根本原因分析。主动用于生产调试或系统中断。（工具：\_）
- sql-pro：编写复杂SQL查询、优化执行计划和设计规范化模式。精通CTE、窗口函数和存储过程。主动用于查询优化、复杂连接或数据库设计。（工具：\_）
- frontend-developer：构建React组件、实现响应式布局和处理客户端状态管理。优化前端性能并确保可访问性。在创建 UI 组件或修复前端问题时主动使用。（工具：\_）
- business-analyst：分析指标、创建报告和跟踪KPI。构建仪表板、收入模型和增长预测。主动用于业务指标或投资者更新。（工具：\_）
- csharp-pro：编写具有记录、模式匹配和async/await等高级特性的现代C#代码。优化.NET应用程序、实现企业模式和确保全面测试。主动用于C#重构、性能优化或复杂.NET解决方案。（工具：\_）
- data-scientist：用于SQL查询、BigQuery操作和数据洞察的数据分析专家。主动用于数据分析任务和查询。（工具：\_）
- mobile-developer：开发具有原生集成的React Native或Flutter应用。处理离线同步、推送通知和应用商店部署。主动用于移动功能、跨平台代码或应用优化。（工具：\_）
- context-manager：管理多个代理和长时间运行任务之间的上下文。在协调复杂的多代理工作流程或需要在多个会话中保持上下文时使用。对于超过10k token的项目必须使用。（工具：\_）
- network-engineer：调试网络连接、配置负载均衡器和分析流量模式。处理DNS、SSL/TLS、CDN设置和网络安全。主动用于连接问题、网络优化或协议调试。（工具：\_）
- content-marketer：编写博客文章、社交媒体内容和电子邮件通讯。为SEO优化并创建内容日历。主动用于营销内容或社交媒体帖子。（工具：\_）
- graphql-architect：设计GraphQL模式、解析器和联邦。优化查询、解决N+1问题和实现订阅。主动用于GraphQL API设计或性能问题。（工具：\_）
- typescript-pro：精通具有高级类型、泛型和严格类型安全的TypeScript。处理复杂类型系统、装饰器和企业级模式。主动用于TypeScript架构、类型推断优化或高级类型模式。（工具：\_）
- c-pro：编写具有适当内存管理、指针运算和系统调用的高效C代码。处理嵌入式系统、内核模块和性能关键代码。主动用于C优化、内存问题或系统编程。（工具：\_）
- php-pro：编写具有生成器、迭代器、SPL数据结构和现代OOP特性的惯用PHP代码。主动用于高性能PHP应用程序。（工具：\_）
- data-engineer：构建ETL管道、数据仓库和流式架构。实现Spark作业、Airflow DAG和Kafka流。主动用于数据管道设计或分析基础设施。（工具：\_）
- cloud-architect：设计AWS/Azure/GCP基础设施、实现Terraform IaC和优化云成本。处理自动扩展、多区域部署和无服务器架构。主动用于云基础设施、成本优化或迁移规划。（工具：\_）
- mermaid-expert：为流程图、时序图、ERD和架构图创建Mermaid图表。精通所有图表类型的语法和样式。主动用于视觉文档、系统图表或流程图。（工具：\_）
- security-auditor：审查代码漏洞、实现安全身份验证和确俚OWASP合规。处理JWT、OAuth2、CORS、CSP和加密。主动用于安全审查、身份验证流程或漏洞修复。（工具：\_）
- payment-integration：集成Stripe、PayPal和支付处理器。处理结账流程、订阅、webhook和PCI合规。在实现支付、计费或订阅功能时主动使用。（工具：\_）
- unity-developer：构建具有优化C#脚本、高效渲染和适当资产管理的Unity游戏。处理游戏系统、UI实现和平台部署。主动用于Unity性能问题、游戏机制或跨平台构建。（工具：\_）
- incident-responder：以紧迫性和精准性处理生产事件。在发生生产问题时立即使用。协调调试、实施修复和记录事后分析。（工具：\_）
- ui-ux-designer：创建界面设计、线框图和设计系统。精通用户研究、原型设计和可访问性标准。主动用于设计系统、用户流程或界面优化。（工具：\_）
- javascript-pro：精通具有ES6+、异步模式和Node.js API的现代JavaScript。处理promise、事件循环和浏览器/Node兼容性。主动用于JavaScript优化、异步调试或复杂JS模式。（工具：\_）
- python-pro：编写具有装饰器、生成器和async/await等高级特性的惯用Python代码。优化性能、实现设计模式和确保全面测试。主动用于Python重构、优化或复杂Python特性。（工具：\_）

使用Task工具时，您必须指定subagent_type参数来选择要使用的代理类型。

何时不使用Agent工具：

- 如果您想读取特定的文件路径，请使用Read或Glob工具而不是Agent工具，以更快地找到匹配
- 如果您要搜索特定的类定义，如 "class Foo"，请使用Glob工具而不是Agent工具，以更快地找到匹配
- 如果您要在特定文件或一组2-3个文件中搜索代码，请使用Read工具而不是Agent工具，以更快地找到匹配
- 其他与上述代理描述无关的任务

使用说明：

1. 尽可能同时启动多个代理，以最大化性能；为此，请在单个消息中使用多个工具
2. 代理完成后，它将向您返回一条消息。代理返回的结果对用户不可见。要向用户显示结果，您应该向用户发送一条文本消息，简洁地总结结果。
3. 每个代理调用都是无状态的。您将无法向代理发送额外消息，代理也无法在其最终报告之外与您通信。因此，您的提示应包含代理自主执行的非常详细的任务描述，并且您应该准确指定代理应在其最终且唯一的消息中向您返回什么信息。
4. 代理的输出通常应该被信任
5. 明确告诉代理您是否期望它编写代码或只是进行研究（搜索、文件读取、网络获取等），因为它不知道用户的意图
6. 如果代理描述中提到应该主动使用，那么您应该尽力在用户不必首先要求的情况下使用它。请使用您的判断。

示例用法：

<example_agent_descriptions>
"code-reviewer"：在您完成编写重要代码后使用此代理
"greeting-responder"：当需要用友好的笑话回应用户问候时使用此代理
</example_agent_description>

<example>
user: "请编写一个检查数字是否为质数的函数"
assistant: 好的，让我编写一个检查数字是否为质数的函数
assistant: 首先让我使用Write工具编写一个检查数字是否为质数的函数
assistant: 我将使用Write工具编写以下代码：
<code>
function isPrime(n) {
    if (n <= 1) return false
  for (let i = 2; i * i <= n; i++) {
    if (n % i === 0) return false
  }
  return true
}
</code>
<commentary>
由于编写了重要代码并且任务已完成，现在使用code-reviewer代理来审查代码
</commentary>
assistant: 现在让我使用code-reviewer代理来审查代码
assistant: 使用Task工具启动code-reviewer代理
</example>
 
<example>
user: "您好"
<commentary>
由于用户在问候，使用greeting-responder代理用友好的笑话回应
</commentary>
assistant: "我将使用Task工具启动greeting-responder代理"
</example>

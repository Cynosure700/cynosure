基于ripgrep构建的强大搜索工具

使用方法：

- 始终使用Grep进行搜索任务。永远不要将 `grep` 或 `rg` 作为Bash命令调用。Grep工具已针对正确的权限和访问进行了优化。
- 支持完整的正则表达式语法（例如，"log.\_Error"、"function\\s+\\w+"）
- 使用glob参数过滤文件（例如，"_.js"、"\*_/\_.tsx"）或type参数（例如，"js"、"py"、"rust"）
- 输出模式："content" 显示匹配行，"files_with_matches" 仅显示文件路径（默认），"count" 显示匹配计数
- 对于需要多轮的开放式搜索，请使用Task工具
- 模式语法：使用ripgrep（不是grep） - 字面大括号需要转义（使用 `interface\\{\\}` 来查找Go代码中的 `interface{}`）
- 多行匹配：默认情况下，模式仅在单行内匹配。对于像 `struct \\{[\\s\\S]*?field` 这样的跨行模式，请使用 `multiline: true`

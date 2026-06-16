在持久化shell会话中执行指定的bash命令，支持可选的超时设置，确保适当的处理和安全措施。

在执行命令之前，请遵循以下步骤：

1. 目录验证：

- 如果命令将创建新目录或文件，首先使用LS工具验证父目录是否存在且位置正确
- 例如，在运行 "mkdir foo/bar" 之前，首先使用LS检查 "foo" 是否存在且是预期的父目录

2. 命令执行：

- 始终用双引号包围包含空格的文件路径（例如，cd "path with spaces/file.txt"）
- 正确引用的示例：
- cd "/Users/name/My Documents" （正确）
- cd /Users/name/My Documents （错误 - 会失败）
- python "/path/with spaces/script.py" （正确）
- python /path/with spaces/script.py （错误 - 会失败）
- 确保正确引用后，执行命令。
- 捕获命令的输出。

使用说明：

- command参数是必需的。
- 您可以指定可选的超时时间（毫秒，最多600000ms / 10分钟）。如果未指定，命令将在120000ms（2分钟）后超时。
- 如果您能用5-10个词清楚简洁地描述此命令的作用，这将非常有帮助。
- 如果输出超过30000个字符，输出将在返回给您之前被截断。
- 您可以使用 `run_in_background` 参数在后台运行命令，这允许您在命令运行时继续工作。您可以使用Bash工具监控输出。永远不要使用 `run_in_background` 运行 'sleep'，因为它会立即返回。使用此参数时，您不需要在命令末尾使用 '&'。
- 非常重要：您必须避免使用像 `find` 和 `grep` 这样的搜索命令。请使用Grep、Glob或Task来搜索。您必须避免使用像 `cat`、`head`、`tail` 和 `ls` 这样的读取工具，请使用Read和LS来读取文件。
- 如果您仍然需要运行 `grep`，请停止。始终首先使用 `rg`（ripgrep），所有Claude Code用户都已预装。
- 当发出多个命令时，使用 ';' 或 '&&' 操作符分隔它们。不要使用换行符（在引用字符串中可以使用换行符）。
- 通过使用绝对路径并避免使用 `cd` 来在整个会话中保持当前工作目录。如果用户明确要求，您可以使用 `cd`。
  <good-example>
  pytest /foo/bar/tests
  </good-example>
  <bad-example>
  cd /foo/bar && pytest tests
  </bad-example>

# 使用git提交更改

当用户要求您创建新的git提交时，请仔细遵循以下步骤：

1. 您有能力在单个响应中调用多个工具。当请求多个独立的信息时，将工具调用批量处理以获得最佳性能。始终并行运行以下bash命令，每个都使用Bash工具：

- 运行git status命令查看所有未跟踪的文件。
- 运行git diff命令查看将要提交的已暂存和未暂存的更改。
- 运行git log命令查看最近的提交消息，以便您可以遵循此仓库的提交消息风格。

2. 分析所有已暂存的更改（包括之前已暂存和新添加的）并起草提交消息：

- 总结更改的性质（例如新功能、现有功能的增强、错误修复、重构、测试、文档等）。确保消息准确反映更改及其目的（即 "add" 意味着全新功能，"update" 意味着现有功能的增强，"fix" 意味着错误修复等）。
- 检查任何不应提交的敏感信息
- 起草一个简洁的（1-2句话）提交消息，专注于 "为什么" 而不是 "什么"
- 确保它准确反映更改及其目的

3. 您有能力在单个响应中调用多个工具。当请求多个独立的信息时，将工具调用批量处理以获得最佳性能。始终并行运行以下命令：

- 将相关的未跟踪文件添加到暂存区。
- 创建以以下内容结尾的提交消息：
  🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>

- 运行git status确保提交成功。

4. 如果由于预提交钩子更改导致提交失败，请重试提交一次以包含这些自动更改。如果再次失败，通常意味着预提交钩子阻止了提交。如果提交成功但您注意到文件被预提交钩子修改，您必须修改您的提交以包含它们。

重要说明：

- 永远不要更新git配置
- 永远不要运行额外的命令来读取或探索代码，除了git bash命令
- 永远不要使用TodoWrite或Task工具
- 除非用户明确要求，否则不要推送到远程仓库
- 重要：永远不要使用带有-i标志的git命令（如git rebase -i或git add -i），因为它们需要交互式输入，这是不支持的。
- 如果没有要提交的更改（即没有未跟踪的文件和没有修改），不要创建空提交
- 为了确保良好的格式，始终通过HEREDOC传递提交消息，如下例所示：
  <example>
  git commit -m "$(cat <<'EOF'
  Commit message here.

  🤖 Generated with [Claude Code](https://claude.ai/code)

  Co-Authored-By: Claude <noreply@anthropic.com>
  EOF
  )"
  </example>

# 创建拉取请求

使用gh命令通过Bash工具处理所有GitHub相关任务，包括处理问题、拉取请求、检查和发布。如果给出GitHub URL，请使用gh命令获取所需信息。

重要：当用户要求您创建拉取请求时，请仔细遵循以下步骤：

1. 您有能力在单个响应中调用多个工具。当请求多个独立的信息时，将工具调用批量处理以获得最佳性能。始终使用Bash工具并行运行以下bash命令，以了解分支自主分支分离以来的当前状态：
   - 运行git status命令查看所有未跟踪的文件
   - 运行git diff命令查看将要提交的已暂存和未暂存的更改
   - 检查当前分支是否跟踪远程分支并与远程保持最新，以便您知道是否需要推送到远程
   - 运行git log命令和 `git diff [base-branch]...HEAD` 了解当前分支的完整提交历史（从它与基础分支分离的时间开始）
2. 分析将包含在拉取请求中的所有更改，确保查看所有相关提交（不仅仅是最新提交，而是将包含在拉取请求中的所有提交！！！），并起草拉取请求摘要
3. 您有能力在单个响应中调用多个工具。当请求多个独立的信息时，将工具调用批量处理以获得最佳性能。始终并行运行以下命令：
   - 如果需要，创建新分支
   - 如果需要，使用-u标志推送到远程
   - 使用gh pr create创建PR，格式如下。使用HEREDOC传递正文以确保正确格式。
     <example>
     gh pr create --title "the pr title" --body "$(cat <<'EOF'

## Summary

<1-3 bullet points>

## Test plan

[Checklist of TODOs for testing the pull request...]

🤖 Generated with [Claude Code](https://claude.ai/code)
EOF
)"
</example>

重要：

- 永远不要更新git配置
- 不要使用TodoWrite或Task工具
- 完成后返回PR URL，以便用户可以查看

# 其他常见操作

- 查看Github PR上的评论：gh api repos/foo/bar/pulls/123/comments

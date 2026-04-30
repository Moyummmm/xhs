# 朱马元

> <span class="icon">&#xe60f;</span> `177 5402 2108`&emsp;&emsp;
> <span class="icon">&#xe7ca;</span> `moyuzo6@qq.com`&emsp;&emsp;
> <span class="icon">&#xe600;</span> [moyuzao](https://github.com/octocat)

<img class="avatar" src="https://avatars.githubusercontent.com/u/583231?v=4">

## &#xe80c; 教育经历

<div class="entry-title">
    <h3>电子科技大学 - 研究生 - 计算机技术</h3> 
    <p>2025.09 - 2028.06</p>
</div>

<div class="entry-title">
    <h3>电子科技大学 - 本科 - 计算机科学与技术</h3> 
    <p>2021.09 - 2025.06</p>
</div>

## &#xe618; 实习经历

<!-- <div alt="entry-title">
    <h3>软件工程师 - 章小鱼有限公司</h3> 
    <p>2008.03 - 2009.07</p>
</div>

作为核心开发成员及技术负责人，主导了八爪生物社交平台（OctoHub）的全栈开发与架构设计。

- 设计并实现独特的"八爪风格"用户交互体系，包括：动态触手消息传递系统、墨水喷溅情感反应功能、自适应伪装个人主页，以促进全球八爪生物和猫之间的社区参与，使用户互动频率提升210%。
- 集成 OAuth 认证，与 GitHub 账户进行同步，为 Octocat 和其他在 GitHub 上活跃的八爪生物提供无缝登录和个人资料同步，将认证流程耗时从12.8s缩短至2.3s，获选GitHub年度最佳身份集成案例。

<div class="entry-title">
    <h3>软件开发实习生 - 八爪科技</h3> 
    <p>2008.06 - 2008.08</p>
</div>

与软件工程师团队合作，使用 Octolang 开发数据可视化仪表盘，为海洋保护工作提供八爪种群趋势的洞察。
- 参与会议和代码审议，按照敏捷章鱼论交付高质量的软件，在紧迫的截止日期内完成任务。
- 协助解决技术问题，展现解决问题的技巧和在快节奏环境下积极主动解决挑战的态度。为项目需求、架构设计和编码标准的文档撰写做出贡献，促进团队成员间的知识共享和新成员的快速适应。 -->

## &#xe635; 项目经历

<div class="entry-title">
    <h3>XiaoHongShu 小红书社区平台</h3>
    <a href="https://github.com/Moyummmm/xhs">github.com/Moyummmm/xhs</a>
</div>

全栈克隆项目，前端 React+TypeScript，后端 Go+Chi 框架，数据库 PostgreSQL+Redis。

- 设计用户体系（注册/登录/关注），基于 JWT 实现无状态认证，支持 token 刷新与主动注销
- 实现笔记 CRUD、点赞、评论功能，采用 Redis 缓存热点数据，接口响应时间 < 50ms
- 点赞去重：基于 Redis Set 实现分布式幂等，避免重复点赞
- 关注列表查询优化：采用 Redis ZSet 存储 Timeline，支撑日均万级查询
- 前端使用 Zustand 管理认证状态，TanStack Query 实现数据缓存与预加载，减少冗余请求 60%+

<div class="entry-title">
    <h3>GitFlix</h3>
    <a href="https://github.com/YiNNx/cmd-wrapped">github.com/octocat/gitflix</a>
</div>

全栈 Web 应用程序，前端使用 Octo.js，后端使用 OctoScript，允许用户发现和评价八爪生物主题电影。
- 实现了一个复杂的推荐算法，分析八爪生物的偏好和观影历史，为八爪生物跨多个流派提供八爪主题的电影推荐，确保了个性化和吸引人的内容发现。
- 使用 JSON Web Tokens 和 bcrypt 实现用户身份验证和授权，用于安全密码哈希。利用 GitHub Actions 进行持续集成和部署，确保流畅高效的开发工作流程。

<div class="entry-title">
    <h3>OctoGithubber</h3>
    <a href="https://github.com/YiNNx/cmd-wrapped">github.com/octocat/gitflix</a>
</div>

一款专门针对八爪生物的 GitHub 活动和贡献的网络应用程序，利用 Octo.js 构建前端，Octolang 构建后端。
- 与 GitHub API 集成，检索和分析八爪生物的存储库统计信息、提交历史和拉取请求活动，提供个性化的见解和可视化，深入了解八爪生物的开源之旅。
- 实现了八爪主题的勋章和成就等游戏化元素，激励和鼓励八爪生物达成编码里程碑，促进持续学习和改进。
- 设计了响应式和直观的仪表板界面，具有八爪主题的数据可视化，使八爪生物能够跟踪进度、设定编码目标，并以有趣和吸引人的方式庆祝成就。

## &#xecfa; 专业技能

- 精通Golang，熟悉Java/C++
- 熟练使用AI工具
- 熟练使用Git、Docker、PostgreSQL等工具
- 通过CET6

<p align="center">
<img src="https://github.com/xiaoxuxiansheng/xtimer/blob/main/common/img/xtimer.png" />
<b>xTimer: go 语言实现的分布式定时器</b>
<br/><br/>
</p>

## 📚 前言
仅以此项目怀念和栋哥、欢哥、康总在前公司共同奋斗的岁月~（还有跟着 devin 学习的日子）<br/><br/>
<img src="https://github.com/xiaoxuxiansheng/xtimer/blob/main/common/img/hezhao.jpeg" height="300px"/>

## 📖 简介
一款依赖于 mysql、redis 组件，基于 go 语言实现的分布式定时器

## 🚀 功能
- 提供定时器 crud 能力
- 基于 cron 表达式定义执行规则
- 支持 http 协议回调下游服务

## 🐧 体验页面
<a href="http://43.143.168.5:5173/login">前端体验页面</a> <br/><br/>
登录账号：test
登录密码：test

## 💡 `xTimer` 技术原理
<a href="https://juejin.cn/post/7174007780104208392">xTimer 实现原理</a> <br/><br/>
<a href="https://juejin.cn/post/7116320697139331103">xTimer 前身 workflow.timer 实现原理</a>

## 🖥 接入 sop
1 用户需要提供好 mysql 和 redis 组件；<br/><br/>
2 在 mysql 中执行 ./common/model/sql 下的建表语句；<br/><br/>
3 ./conf.yml 中填写 mysql dsn 以及 redis 账号密码；<br/><br/>
4 运行 main 函数. <br/><br/>

## 📊 定制参数
在 ./conf.yml 中暴露了更多的参数供用户自主选择定义，这部分均配置默认值兜底,用户如需修改,请先了解技术原理,在熟悉源代码的 基础上可自行调节


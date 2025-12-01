//开发一个个人博客系统的后端
//采用go+gin+gorm+mysql(5.7)+jwt
//本地数据库
    mysql 连接信息：
        主机：172.18.112.82
        端口：3306
        数据库名：blog
        用户名：root
        密码：root
//功能  
使用 在04blog目录下使用 go mod init 初始化项目依赖管理。
//1.用户注册登录（包含 JWT 认证）
    users 表：存储用户信息，包括 id 、 username 、 password 、 email 等字段。
//2.发布博客
    posts 表：存储博客文章信息，包括 id 、 title 、 content 、 user_id （关联 users 表的 id ）、 created_at 、 updated_at 等字段。
    comments 表：存储文章评论信息，包括 id 、 content 、 user_id （关联 users 表的 id ）、 post_id （关联 posts 表的 id ）、 created_at 等字段。

//4.点赞博客
    likes 表：存储用户点赞信息，包括 id 、 user_id （关联 users 表的 id ）、 post_id （关联 posts 表的 id ）、 created_at 等字段。
//5.评论博客
    comments 表：存储文章评论信息，包括 id 、 content 、 user_id （关联 users 表的 id ）、 post_id （关联 posts 表的 id ）、 created_at 等字段。

   
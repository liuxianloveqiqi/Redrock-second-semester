# 基于Ubuntu 18.04作为操作系统镜像
FROM ubuntu:18.04

# 安装MySQL
RUN apt-get update && apt-get install -y mysql-server

# 复制MySQL配置文件到容器中
COPY my.cnf /etc/mysql/my.cnf

# 定义MySQL的数据目录
VOLUME /var/lib/mysql

# 定义MySQL的默认端口
EXPOSE 3306

# 设置容器启动时执行的命令
CMD ["mysqld"]

# 构建镜像时运行的命令
RUN service mysql start && mysql -u root -e "CREATE DATABASE mydb"

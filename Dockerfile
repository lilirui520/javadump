# 使用官方的 OpenJDK 镜像作为基础镜像
FROM openjdk:11-jdk-slim

# 设置工作目录
WORKDIR /app

# 将 OOMSimulator.java 文件复制到容器中
COPY OOMSimulator.java .

# 编译 Java 文件
RUN javac OOMSimulator.java

# 设置 JVM 参数以触发 OOM 并生成堆转储文件
ENV JAVA_OPTS="-Xmx10m -XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=/app/dump"

# 创建堆转储文件目录
RUN mkdir -p /app/dump

# 运行程序
CMD java ${JAVA_OPTS} OOMSimulator

###docker run -v $(pwd)/dump:/app/dump oom-simulator
FROM k8s-harbor.exchange.local/base-images/openjdk:8u342-jdk

ARG APP_NAME # 项目名称，比如APP_NAME="account-provider"
ARG JAVA_OPTS # java启动参数，比如JAVA_OPTS="-Xmx1024m -Xms1024m"
ARG TIMEZONE

#ENV TINI_VERSION v0.19.0
#ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /bin/tini
ADD tini /bin/tini
RUN chmod +x /bin/tini
ENTRYPOINT ["/bin/tini", "--"]

RUN groupadd app -g 2000 && useradd app -u 2000 -g 2000 && \
    mkdir -p /usr/local/apps && mkdir -p /usr/local/logs && \
    mkdir -p /data && mkdir -p /home/app && \
    mkdir -p /usr/local/apps/${APP_NAME} && \
    echo "alias egrep='egrep --color=auto'\nalias grep='grep --color=auto'\nalias ll='ls -l --color=auto'\nalias ls='ls --color=auto'"> /home/app/.bashrc 

COPY ${APP_NAME}.jar /usr/local/apps/

RUN chown -R app:app /usr/local/apps /usr/local/logs /data /home/app

ENV JAVA_OPTS=$JAVA_OPTS APP_NAME=$APP_NAME TZ=$TIMEZONE 

WORKDIR /usr/local/apps/

USER app

CMD ["/bin/sh", "-c", "/usr/local/openjdk-8/bin/java $JAVA_OPTS -Ddubbo.reference.check=false -Denv=pro -Dpro_meta=http://config.apollo.internal.com:8080 -jar ${APP_NAME}.jar"]

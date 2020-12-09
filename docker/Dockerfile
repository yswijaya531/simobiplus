FROM devops-registry.smartfren.com/devops/alpine_net_tz:3.11.3

LABEL maintainer="devops@smartfren.com"

COPY simobiplus /opt/
COPY simobi.key /opt/
COPY simobi.key.pub /opt/
COPY simobi-ssl.crt /opt/
COPY simobi-ssl.key /opt/
COPY backend.key.pub /opt/
WORKDIR /opt/

RUN export http_proxy=http://10.0.0.178:8080; \
    export https_proxy=http://10.0.0.178:8080; \
    apk add tzdata
ENV TZ Asia/Jakarta

#RUN echo '10.39.8.6 api.banksinarmas.com' >> /etc/hosts

RUN [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf
#EXPOSE 9661
CMD ["/opt/simobiplus"]
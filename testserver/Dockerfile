FROM		ubuntu:14.04
MAINTAINER	Guillaume J. Charmes <guillaume@charmes.net>

RUN		apt-get update
RUN		apt-get install -y unzip socat xvfb gsettings-desktop-schemas openjdk-7-jre && rm -rf /var/lib/apt/lists/*
ENV		JAVA_HOME /usr/lib/jvm/java-7-openjdk-amd64

EXPOSE		4002
EXPOSE		4003

ENV		IB_LOGIN	fdemo
ENV		IB_PASSWORD	demouser

ADD		.	  /src

CMD		cd /src && ./ibgwdocker

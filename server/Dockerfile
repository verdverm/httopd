# FROM google/debian:wheezy
FROM dockerfile/ubuntu:latest
MAINTAINER verdverm@gmail.com

# Update stuff
RUN apt-get update
RUN apt-get upgrade -y

# Install Nginx.
RUN \
  add-apt-repository -y ppa:nginx/stable && \
  apt-get update && \
  apt-get --no-install-recommends install -y nginx && \
  echo "\ndaemon off;" >> /etc/nginx/nginx.conf && \
  chown -R www-data:www-data /var/lib/nginx

# Define mountable directories.
VOLUME ["/data", "/etc/nginx/sites-enabled", "/var/log/nginx"]

# Install Python Setuptools
RUN apt-get --no-install-recommends install -y python-setuptools build-essential python-dev

# Install pip
RUN easy_install pip

ADD . /httopd
WORKDIR /httopd

# Install requirements.txt
RUN pip install -r requirements.txt

# Expose ports.
EXPOSE 5000 8080

CMD /httopd/run.sh

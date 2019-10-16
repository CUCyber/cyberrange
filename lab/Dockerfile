# Flag 10/11
FROM golang:alpine AS build-flag-ten
RUN apk --no-cache add build-base git bzr mercurial gcc
ADD ./flag_ten /src
RUN cd /src && go build -o flag_ten

# Flag 11
FROM golang:alpine AS build-flag-eleven
RUN apk --no-cache add build-base git bzr mercurial gcc
ADD ./flag_eleven /src
RUN cd /src && go build -o flag_eleven

# Flag 12
FROM golang:alpine AS build-flag-twelve
RUN apk --no-cache add build-base git bzr mercurial gcc
ADD ./flag_twelve /src
RUN cd /src && go build -o flag_twelve



FROM ubuntu:16.04

RUN apt-get update && apt-get install -y openssh-server zip unzip
RUN mkdir /var/run/sshd
RUN echo 'root:abc123***' | chpasswd
RUN sed -i 's/PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config

# SSH login fix. Otherwise user is kicked off after login
RUN sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd

ENV NOTVISIBLE "in users profile"
RUN echo "export VISIBLE=now" >> /etc/profile

# Teaching setup
## First User
RUN useradd -ms /bin/bash one && echo "one:newpassword" | chpasswd
USER one
WORKDIR /home/one

RUN touch flag{1}

RUN mkdir dir_one
RUN touch dir_one/flag{2}

RUN mkdir "dir two"
RUN touch dir\ two/flag{3}

RUN echo "flag{5}" >> flag_five

RUN touch .flag{6}

RUN echo "flag{7}" >> .flag_seven
RUN chmod -r .flag_seven

## Second User
USER root
RUN useradd -ms /bin/bash two && echo 'two:newpassword' | chpasswd
USER two
WORKDIR /home/two

RUN touch flag{8}


## Root User
USER root
WORKDIR /root

RUN touch flag{9}


## Third User
RUN useradd -ms /bin/bash three && echo 'three:newpassword' | chpasswd

COPY --from=build-flag-ten /src/flag_ten /home/three/flag_ten
RUN chown three:three /home/three/flag_ten
COPY --from=build-flag-eleven /src/flag_eleven /home/three/flag_eleven
RUN chown three:three /home/three/flag_eleven && chmod -x /home/three/flag_eleven
COPY --from=build-flag-twelve /src/flag_twelve /home/three
RUN chown three:three /home/three/flag_twelve


## Fourth User
RUN useradd -ms /bin/bash four && echo 'four:newpassword' | chpasswd && echo "flag{13}" | base64 > /home/four/flag_thirteen
COPY ./flag_fourteen /home/four
COPY ./flag_fifteen.zip /home/four
COPY ./flag_sixteen.tar /home/four

EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]

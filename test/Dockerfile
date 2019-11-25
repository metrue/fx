FROM ubuntu:16.04

RUN apt-get update && apt-get install -y openssh-server curl
RUN mkdir -p ~/.ssh
RUN echo 'ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDPVROluD9aW8YEsHiMefr0Yk70TzMJ+yRXkTN0DSDQje6fycffZaxI4vb5JO/tfXkTQCg+uo3t9YVQU3ceFAPpnznCnCr3jnOo7s2BqV5zDRjIW/fG3MLuVyZKvecA5RDIj2WLfvlsev+J6LI/Q/kMr9i8dI9BHp5B3u8Nv3sePEzKU9YRnTd/UTbSdAHKqfpGhgwZEI00q3iiP6f5DKVXZ4b7ZVEsV3cPVrRskurYClSMd32/yaJ+68mFlpwTKI/aq7tZBd5lLsAsd2IxshGE23g4bU04GeeJ76tFT7BvDyL8woshECisRHSdEsdlY9MXIcC/a4hIV4baHXJDkFrf minghe@oldmac.local' >> ~/.ssh/authorized_keys
RUN mkdir /var/run/sshd
RUN echo 'root:THEPASSWORDYOUCREATED' | chpasswd
RUN sed -i 's/PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config

# SSH login fix. Otherwise user is kicked off after login
RUN sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd

ENV NOTVISIBLE "in users profile"
RUN echo "export VISIBLE=now" >> /etc/profile

EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]

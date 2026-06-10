FROM alpine:latest
RUN apk add --no-cache openssh-server rsync
RUN ssh-keygen -A
RUN adduser -D -s /bin/sh testuser && \
    echo "testuser:testpass" | chpasswd && \
    mkdir -p /home/testuser/.ssh && \
    chmod 700 /home/testuser/.ssh
COPY id_rsa.pub /home/testuser/.ssh/authorized_keys
RUN chmod 600 /home/testuser/.ssh/authorized_keys && \
    chown -R testuser:testuser /home/testuser/.ssh
RUN echo "PermitRootLogin no" >> /etc/ssh/sshd_config && \
    echo "PasswordAuthentication yes" >> /etc/ssh/sshd_config && \
    echo "PubkeyAuthentication yes" >> /etc/ssh/sshd_config
EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]

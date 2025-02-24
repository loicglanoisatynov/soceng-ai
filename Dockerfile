# FROM golang:1.18-alpine
FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o ./bin/main ./cmd/main/main.go

RUN go build -o ./bin/server ./cmd/server/main.go

# RUN echo "[Unit]" >> /etc/systemd/system/server.service
# RUN echo "Description=Serveur d'ingénierie sociale" >> /etc/systemd/system/server.service
# RUN echo "After=network.target" >> /etc/systemd/system/server.service
# RUN echo "" >> /etc/systemd/system/server.service
# RUN echo "[Service]" >> /etc/systemd/system/server.service
# RUN echo "ExecStart=/home/user/server start" >> /etc/systemd/system/server.service
# RUN echo "ExecStop=/home/user/server stop" >> /etc/systemd/system/server.service
# RUN echo "Restart=always" >> /etc/systemd/system/server.service
# RUN echo "User=user" >> /etc/systemd/system/server.service
# RUN echo "Group=user" >> /etc/systemd/system/server.service
# RUN echo "WorkingDirectory=/home/user" >> /etc/systemd/system/server.service
# RUN echo "" >> /etc/systemd/system/server.service
# RUN echo "[Install]" >> /etc/systemd/system/server.service
# RUN echo "WantedBy=multi-user.target" >> /etc/systemd/system/server.service

# RUN echo "root:root" | chpasswd

# RUN sudo systemctl daemon-reload
# RUN sudo systemctl enable server
# RUN sudo systemctl start server


EXPOSE 80

# Execute the server
CMD ["/bin/bash"]
all : build stop remove start
.PHONY : all

build :
	sudo docker build -t google-calendar-server .

stop:
	sudo docker stop google-calendar-server

remove :
	sudo docker rm google-calendar-server

start :
	sudo docker run --name google-calendar-server -p 9000:9000 -d --restart=always google-calendar-server
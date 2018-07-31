deploy: clean build push

clean:
	cd website && rm -fR build

run:
	cd website && bundle exec middleman server

build:
	cd website && bundle exec middleman build

push:
	firebase deploy

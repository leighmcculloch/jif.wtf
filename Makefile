deploy: clean build push

clean:
	rm -fR website/build

build: build-website

build-website:
	[ -t 0 ] || (cd website && bundle install)
	cd website && bundle exec middleman build

push: push-service push-website

push-service:
	gcloud app deploy --project jif-wtf-bf47b service/app.yaml

push-website: .bin/node_modules/.bin/firebase
	[ ! -t 0 ] || .bin/node_modules/.bin/firebase login --no-localhost
	.bin/node_modules/.bin/firebase deploy

.bin:
	mkdir .bin

.bin/node_modules/.bin/firebase: .bin
	yarn --no-lockfile --modules-folder=.bin/node_modules add firebase-tools

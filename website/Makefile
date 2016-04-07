export CLOUDFLARE_ZONE = f24e7be338605506f4cd84f8ff6dda42

deploy: clean build push cdn

clean:
	rm -fR build

build:
	bundle exec middleman build

push:
	bundle exec middleman s3_sync --force

cdn:
	curl -X DELETE "https://api.cloudflare.com/client/v4/zones/$(CLOUDFLARE_ZONE)/purge_cache" \
		-H "X-Auth-Email: $(CLOUDFLARE_EMAIL)" \
		-H "X-Auth-Key: $(CLOUDFLARE_CLIENT_API_KEY)" \
		-H "Content-Type: application/json" \
		--data '{"purge_everything":true}'


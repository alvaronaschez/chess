fe: fe-vue

be:
	go run .

fe-vue:
	(cd fe-chess && npm run dev)

fe-web-components:
	npx live-server

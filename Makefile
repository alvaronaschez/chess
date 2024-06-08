fe: fe-vue

be:
	(cd api && go run .)

fe-vue:
	(cd fe-chess && npm run dev)

fe-web-components:
	(cd fe-web-components && npx live-server)

fe: fe-vue

be:
	(cd backend && go run .)

fe-vue:
	(cd frontend && npm run dev)

fe-web-components:
	(cd frontend-web-components && npx live-server)

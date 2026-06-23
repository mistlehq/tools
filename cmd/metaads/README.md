# Meta Ads CLI

`metaads` is a thin command-line and MCP wrapper around Meta Graph API / Marketing API.

It does not mint, store, refresh, or inspect Meta credentials. Set `METAADS_GRAPH_BASE_URL` to the versioned Graph API base URL or a local proxy that injects auth:

```sh
export METAADS_GRAPH_BASE_URL="https://graph.facebook.com/v25.0"
```

## Commands

```sh
metaads help
metaads version
metaads auth test
metaads graph request --method GET --path /me
metaads graph request --method GET --path /act_<account-id>/campaigns --params '{"fields":"id,name"}'
metaads ad-accounts list
metaads ad-accounts get --id act_<account-id>
metaads campaigns search|get|create|update|delete
metaads adsets search|get|create|update|delete
metaads ads search|get|create|update|delete
metaads creatives search|get|create
metaads insights get --id act_<account-id>
metaads targeting search --type adinterest --query running
metaads mcp serve
```

`metaads graph request` is the complete API coverage surface. Named commands are convenience wrappers for progressive discovery of common Marketing API workflows.

## Auth

Auth is expected to be handled outside this binary. In Mistle, egress injects the Meta access token. For local direct usage, run through a proxy or provide auth at the network layer.

Meta permissions and ad-account access are managed by the user in Meta Business Manager / App settings. Provider authorization errors are returned as Meta returns them.

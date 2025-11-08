package sdk

//go:generate sdk-generator --level info --language typescript --output js  ../...
//go:generate pnpm --dir js install
//go:generate pnpm --dir js format

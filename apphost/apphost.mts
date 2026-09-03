import { createBuilder } from './.aspire/modules/aspire.mjs';

const builder = await createBuilder();

await builder.addGoApp('gateway', '../services/gateway', {
  packagePath: './cmd/gateway'
}).withHttpEndpoint({env: 'PORT'});

await builder.build().run();

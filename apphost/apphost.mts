import { createBuilder } from './.aspire/modules/aspire.mjs';

const builder = await createBuilder();

var gitRepo = await builder.addGoApp('git-repo', '../services/git-repo', {
  packagePath: './cmd/git-repo'
}).withHttpEndpoint({ name: 'grpc', env: 'GRPC_PORT'}).asHttp2Service();

await builder.addGoApp('gateway', '../services/gateway', {
  packagePath: './cmd/gateway',
}).withHttpEndpoint({env: 'PORT'})
  .withReference(gitRepo);


await builder.build().run();

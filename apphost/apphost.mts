import { createBuilder, EndpointProperty } from './.aspire/modules/aspire.mjs';

const builder = await createBuilder();

var gitRepo = await builder.addGoApp('git-repo', '../services/git-repo', {
  packagePath: './cmd/git-repo'
}).withHttpsEndpoint({ name: 'grpc', env: 'GRPC_PORT'}).asHttp2Service();

await builder.addGoApp('gateway', '../services/gateway', {
  packagePath: './cmd/gateway',
}).withHttpEndpoint({env: 'PORT'})
  .withReference(gitRepo) // keeps dependency ordering
  .withEnvironment(
    'GIT_REPO_GRPC',
    gitRepo.getEndpoint('grpc').property(EndpointProperty.HostAndPort)
  );


await builder.build().run();

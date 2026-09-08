import path from 'path';
import { createBuilder, EndpointProperty, ProtocolType } from './.aspire/modules/aspire.mjs';

const builder = await createBuilder();

const sshKey = path.resolve('../id_rsa');

await builder.addGoApp('git-ssh', '../services/ssh', {
  packagePath: './cmd/ssh'
})
.withEnvironment("SSH_HOST_KEY", sshKey)
.withEndpoint({
    name: 'ssh',
    port: 2222,
    env: 'SSH_LISTEN_PORT',
    protocol: ProtocolType.Tcp,
});

var gitRepo = await builder.addGoApp('git-repo', '../services/git-repo', {
  packagePath: './cmd/git-repo'
})
.withHttpsEndpoint({ name: 'grpc', env: 'GRPC_PORT'})
.asHttp2Service();


await builder.addGoApp('gateway', '../services/gateway', {
  packagePath: './cmd/gateway',
})
.withHttpEndpoint({port: 8080, env: 'PORT'})
.withReference(gitRepo) // keeps dependency ordering
.withEnvironment(
  'GIT_REPO_GRPC',
  gitRepo.getEndpoint('grpc').property(EndpointProperty.HostAndPort)
);

await builder.build().run();

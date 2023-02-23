import * as childProcess from 'child_process';
import * as process from 'process';

const VERSION = '0.0.1';

function chooseBinary() {
  const platform = process.platform;
  const arch = process.arch;

  // Currently only supports Linux on X86_64
  if (platform === 'linux' && arch === 'x64') {
    return `action-linux-amd64-${VERSION}`;
  }

  // If we're here the environment is incorrect
  console.error(`Unsupported platform (${platform}) and architecture (${arch})`);
  return '';
}

function main() {
  // Choose correct Go binary
  const binary = chooseBinary();
  if (binary === '') {
    process.exit(1);
  }
  const actionBin = path.join(`${__dirname}`, 'bin', `${binary}`)
  // Execute main Go program
  const spawnSyncReturns = childProcess.spawnSync(actionBin, {
    stdio: 'inherit'
  });
  // Check exit code
  const status = spawnSyncReturns.status;
  if (typeof status === 'number') {
    process.exit(status);
  }
  // If we're here, something went wrong
  process.exit(1);
}

if (require.main === module) {
  main();
}

import { spawnSync } from 'child_process';
import path from 'path';
import process from 'process';

const VERSION = '0.0.1';

main();

function main() {
  console.log("Node Version: ", process.version);

  // Execute the main Go program
  const spawnSyncReturns = spawnSync('go', ['run', 'main.go'], {
    cwd: path.join(process.cwd(), 'action'),
    env: process.env,
    stdio: 'inherit',
    encoding: 'utf-8'
  });

  // Check the exit code
  const status = spawnSyncReturns.status;
  if (typeof status === 'number') {
    process.exit(status);
  }

  // If we're here, something went wrong
  process.exit(1);
}


import { spawnSync } from 'child_process';
import path from 'path';
import process from 'process';

const VERSION = '0.0.1';

main();

function main() {
  console.log("Node Version: ", process.version);

  // Compile Go program
  console.log("Building ...");
  const goBuildReturns = spawnSync('go', ['build', '-o', '../action.bin'], {
    cwd: path.join(process.cwd(), 'action'),
    stdio: 'inherit',
    encoding: 'utf-8'
  });
  // Check the exit code
  const statusBuild = goBuildReturns.status;
  if (statusBuild !== 0) {
    process.exit(statusBuild);
  }

  // Run Go program
  console.log("Running ...");
  const binRunReturns = spawnSync('./action.bin', [], {
    cwd: process.cwd(),
    stdio: 'inherit',
    encoding: 'utf-8'
  });
  // Check the exit code
  const statusRun = binRunReturns.status;
  if (typeof statusRun === 'number') {
    process.exit(statusRun);
  }

  // If we're here, something went wrong
  process.exit(1);
}


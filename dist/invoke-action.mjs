import { spawnSync } from 'child_process';
import path from 'path';
import process from 'process';

const VERSION = '0.3.0';

main();

function main() {
  console.log("Node Version: ", process.version);

  // Compile Go program
  console.log("Building ...");
  var actionGoDir = path.join( path.dirname(process.argv[1]), '..', 'action');

  const goBuildReturns = spawnSync('go', ['build', '-o', path.join(process.cwd(), 'action.bin')], {
    cwd: actionGoDir,
    stdio: 'inherit',
    encoding: 'utf-8'
  });
  // Check the exit code
  const statusBuild = goBuildReturns.status;
  if (statusBuild !== 0) {
    console.log("Building failed");
    process.exit(statusBuild);
  }
  console.log("Building OK");

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
    if (statusRun !== 0) {
      console.log("ERROR: Running failed!");
    }
    process.exit(statusRun);
  }

  // If we're here, something went wrong
  process.exit(1);
}


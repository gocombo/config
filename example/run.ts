#!/usr/bin/env ts-node

import { Writable } from 'stream';
import { ConfigValue } from '../src/value';
import loadConfig from './config/load';
import { HelloConfig } from './config/types';

interface HelloProducer {
  produceHelloMessage() : string
}

function createHelloProducer(cfg: HelloConfig): HelloProducer {
  // Each config value is accessed via
  return {
    produceHelloMessage() : string {
      // Config value is accessed via function call
      // this allows a dynamic (runtime) config updates
      // and should be a preferable method to access config values
      const message = cfg.message();

      return message;
    },
  };
}

interface SayHelloService {
  writeHello(out: Writable): void
}

function createSayHelloService(
  helloProducer: HelloProducer,
  sayHelloTimes: ConfigValue<number>,
): SayHelloService {
  return {
    writeHello(out: Writable) {
      const times = sayHelloTimes();
      out.write(`Producing hello message ${times} times\n`);
      for (let i = 0; i < times; i += 1) {
        out.write(helloProducer.produceHelloMessage());
        out.write('\n');
      }
      out.write('Done\n');
    },
  };
}

async function main() {
  // Load config just once on app boot time
  // Then pass config values/structures to consumer components
  // Avoid sharing toplevel config reference with consumer components
  const { config, waitPendingSources } = loadConfig();

  // The waitPendingSources is a promise that allows waiting for async sources to load
  // or refresh. For long running processes it may be done just once on app start.
  // For cloud functions it has to be done per function invocation.
  // An example can be found in this file: https://github.com/freshcutgg/token-transactions/blob/dev/index.ts
  // Note: It will only fetch it once or if new values needs to be refreshed, but not each time
  // the function is invoked.
  await waitPendingSources;

  const helloProducer = createHelloProducer(config.hello);
  const { writeHello } = createSayHelloService(helloProducer, config.sayHelloTimes);

  writeHello(process.stdout);
}

main();

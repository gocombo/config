import path from 'path';
import {
  createConfigLoader,
  LoadConfigResult,
  createJSONDataSource,
  createKeyValueDataSource,
} from '../..';
import { Config } from './types';

// GCP_PROJECT env variable to define files to load configs from
export default function loadConfig(
  configFileName = `${process.env.GCP_PROJECT || 'local'}`,
) : LoadConfigResult<Config> {
  const loader = createConfigLoader<Config>({
    dataSources: [
      // Order of data sources makes a difference
      // Last one has a priority and will "override" values
      // of a previous one, if such values are available in a source.

      // default.json defines base config with initial (default) values
      createJSONDataSource(path.join(__dirname, 'default.json')),

      // environment specific file allows overriding defaults
      createJSONDataSource(path.join(__dirname, `${configFileName}.json`)),

      // <user> configs can be used to let devs override values locally without committing
      createJSONDataSource(path.join(__dirname, `${configFileName}-user.json`), {
        ignoreMissingFile: true,
      }),

      // Allow overriding some values via environment variables
      createKeyValueDataSource({
        'hello/message': process.env.HELLO_MESSAGE,
      }),
    ],
  });
  return loader((defineValue) => ({
    hello: {
      message: defineValue({ path: 'hello/message' }),
    },
    sayHelloTimes: defineValue({ path: 'sayHelloTimes' }),
  }));
}

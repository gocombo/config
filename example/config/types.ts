import { ConfigValue } from '../..';

export interface HelloConfig {
  message: ConfigValue<string>
}

export interface Config {
  sayHelloTimes: ConfigValue<number>
  hello: HelloConfig
}

import {sendWithPromise} from 'chrome://resources/js/cr.m.js';

export interface HNSInternalsProxy {
  initialize(): Promise<string>;
}

export class HNSInternalsProxyImpl implements HNSInternalsProxy {
  initialize(): Promise<string> {
    return sendWithPromise('initialize');
  }

  static getInstance(): HNSInternalsProxy {
    return instance || (instance = new HNSInternalsProxyImpl());
  }

  static setInstance(obj: HNSInternalsProxy) {
    instance = obj;
  }
}

let instance: HNSInternalsProxy|null = null;

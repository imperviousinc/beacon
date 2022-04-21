// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

import {sendWithPromise} from 'chrome://resources/js/cr.m.js';

export interface WelcomeProxy {
  initialize(): Promise<string>;
}

export class WelcomeProxyImpl implements WelcomeProxy {
  initialize(): Promise<string> {
    return sendWithPromise('initialize');
  }

  static getInstance(): WelcomeProxy {
    return instance || (instance = new WelcomeProxyImpl());
  }

  static setInstance(obj: WelcomeProxy) {
    instance = obj;
  }
}

let instance: WelcomeProxy|null = null;

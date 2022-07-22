import { Fetcher } from 'openapi-typescript-fetch';
import type { paths } from './proto/descan';
import type { definitions } from '../lib/proto/descan';

import { writable, type Writable } from "svelte/store";

export function CreateClient(baseUrl: string) {
  const fetcher = Fetcher.for<paths>();
  fetcher.configure({
    baseUrl: baseUrl
  });
  return fetcher;
}


export class DescanClient {
  fetcher;
  endpoint: string;
  alive: Writable<boolean>;
  message: Writable<string>;
  network: number;

  constructor(endpoint: string, network: number) {
    const fetcher = Fetcher.for<paths>();
    fetcher.configure({
      baseUrl: endpoint
    });
    this.fetcher = fetcher;
    this.endpoint = endpoint;
    this.alive = writable(false);
    this.message = writable("");
    this.network = network;
    setInterval(this.healthCheck.bind(this), 5000);
    this.healthCheck();
  }

  setNetwork(network: number) {
    this.network = network;
  }

  async healthCheck() {
    if (this.network == 0) {
      return;
    }
    try {
      const stat = await this.getStatus(this.network.toString());
      this.alive.set(true);
      if (stat.message) {
        this.message.set(stat.message);
      }
    } catch (err) {
      this.message.set("");
      this.alive.set(false);
      console.log(err);
    }
  }

  setEndpoint(endpoint: string) {
    const fetcher = Fetcher.for<paths>();
    fetcher.configure({
      baseUrl: endpoint
    });
    this.fetcher = fetcher;
    this.endpoint = endpoint;
    this.healthCheck();
  }

  getEndpoint(): string {
    return this.endpoint;
  }

  async getStatus(chainId: string): Promise<definitions["descanHello"]> {
    const hello = this.fetcher.path('/v1/descan/hello').method('post').create();
    const response = await hello({ message: 'descan r0', chainId: chainId });
    return response.data;
  }

  async getTxsForAddress(address: string, page: string, offset: string): Promise<definitions['descanTxByAddressResponse'] | undefined> {
    const fetchTxs = this.fetcher.path('/v1/descan/tx_by_address').method('post').create();
    const response = await fetchTxs({
      address: address,
      offset: offset,
      page: page
    });
    return response.data;
  }
}

import { Fetcher } from 'openapi-typescript-fetch'
import type { paths } from './proto/descan'

export function CreateClient(baseUrl:string) {
  const fetcher = Fetcher.for<paths>();
  fetcher.configure({
    baseUrl: baseUrl
  });
  return fetcher;
}

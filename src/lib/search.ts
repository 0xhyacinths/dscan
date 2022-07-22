export enum ResultType {
  Address,
  Block,
  Transaction
}

export interface SearchResult {
  query: string,
  type: ResultType,
  page: number | undefined,
}

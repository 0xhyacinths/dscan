<script lang="ts">
	import { Col, Row, Alert, Spinner, Card, CardBody, CardHeader, Table } from 'sveltestrap';
	import { ethers } from 'ethers';
	import { onMount, afterUpdate, createEventDispatcher } from 'svelte';
	import { State } from '../lib/state';
	// @ts-ignore
	import borc from 'borc';

	import type { DescanClient } from '../lib/api';
	import type { definitions } from '../lib/proto/descan';
	const dispatch = createEventDispatcher();
	import type { SearchResult } from '../lib/search';
	import { ResultType } from '../lib/search';

	function formatEth(num: ethers.BigNumber | null | undefined): string {
		return ethers.utils.formatUnits(num!, 'ether');
	}

	export let query: string;
	export let client: DescanClient;
	let lastQuery: string;

	let state: State = State.Loading;
	let err = '';

	let balance: ethers.BigNumber;
	let contractCode: string;
	let totalTx: number;
	let ensAddress: string | null;
	let cborData: string | null;

	let txs: definitions['descanTransactionSummary'][] | undefined;

	onMount(async () => {
		lastQuery = query;
		await fetchData();
	});

	afterUpdate(async () => {
		if (lastQuery == query) {
			return;
		}
		lastQuery = query;
		await fetchData();
	});

	async function fetchData() {
    state = State.Loading;
		if (window.ethereum === undefined) {
			state = State.Error;
			err = 'No Ethereum provider found.';
			return;
		}
		const provider = new ethers.providers.Web3Provider(window.ethereum, 'any');
		try {
			balance = await provider.getBalance(query);
			contractCode = await provider.getCode(query);
			if (contractCode.length > 2) {
				cborData = parseCBOR(contractCode);
			} else {
				cborData = null;
			}
			totalTx = await provider.getTransactionCount(query);
			ensAddress = await provider.lookupAddress(query);
		} catch (e) {
			err = (e as any).toString();
			state = State.Error;
			return;
		}
		try {
			txs = await client.getTxsForAddress(query);
		} catch (e) {
			console.log('no server', e);
		}
		state = State.Loaded;
	}

	function parseCBOR(code: string): string | null {
		const length = ethers.BigNumber.from(`0x${code.substring(code.length - 4)}`).toNumber();
		const cborData = code.substring(code.length - (4 + length * 2), code.length - 4);
		if (cborData) {
			try {
				const data = borc.decode(cborData);
				let storage: any = {};
				for (const property in data) {
					storage[property] = buf2hex(data[property] as Uint8Array);
				}
				return JSON.stringify(storage, null, 2);
			} catch (error) {
				return null;
			}
		}
		return null;
	}

	function buf2hex(buffer: Uint8Array): string {
		return [...new Uint8Array(buffer)].map((x) => x.toString(16).padStart(2, '0')).join('');
	}

	function linkAddress(address: string | undefined) {
		if (address) {
			dispatch('search', {
				query: address,
				type: ResultType.Address
			} as SearchResult);
		}
	}

	function linkBlock(block: string | undefined) {
		if (block) {
			dispatch('search', {
				query: block,
				type: ResultType.Block
			} as SearchResult);
		}
	}

	function linkTx(tx: string | undefined) {
		if (tx) {
			dispatch('search', {
				query: tx,
				type: ResultType.Transaction
			} as SearchResult);
		}
	}

	function isCurrent(address: string | undefined): boolean {
		if (address) {
			return address.toLowerCase() == query.toLowerCase();
		}
		return false;
	}

	function abbr(value: string | undefined): string {
		const partLength = 6;
		if (value) {
			return `${value.slice(0, partLength + 2)}...${value.substring(value.length - partLength)}`;
		}
		return '';
	}
</script>

<Row>
	<Col
		>{#if state == State.Error}
			<Alert color="danger">
				<h4 class="alert-heading text-capitalize">Error loading address</h4>
				{err}
			</Alert>
		{:else if state == State.Loading}
			<div class="text-center">
				<Spinner color="dark" />
			</div>
		{:else if state == State.Loaded}
			<Card class="mb-3">
				<CardHeader>
					Address {query}
				</CardHeader>
				<CardBody>
					<Table responsive borderless class="table-nomargin">
						<tbody>
							<tr>
								<th class="titleWidth" scope="row">Balance</th>
								<td>{formatEth(balance)} Ether</td>
							</tr>
							<tr>
								<th class="titleWidth" scope="row">Total Transactions</th>
								<td>{totalTx}</td>
							</tr>
							{#if ensAddress}
								<tr>
									<th class="titleWidth" scope="row">ENS Address</th>
									<td>{ensAddress}</td>
								</tr>
							{/if}
							{#if contractCode.length > 2}
								<tr>
									<th class="titleWidth" scope="row">Contract code</th>
									<td
										><textarea
											class="form-control contractCode"
											rows="10"
											name="text"
											readonly
											bind:value={contractCode}
										/></td
									>
								</tr>
							{/if}
							{#if cborData}
								<tr>
									<th class="titleWidth" scope="row">Metadata</th>
									<td
										><textarea
											class="form-control contractCode"
											rows="4"
											name="text"
											readonly
											bind:value={cborData}
										/></td
									>
								</tr>
							{/if}
						</tbody>
					</Table>
					{#if txs}
						<hr />
						<Table responsive borderless class="table-nomargin">
							<tbody>
								<tr>
									<th>Tx Hash</th>
									<th>Block</th>
									<th>From</th>
									<th>To</th>
									<th>Value (Ether)</th>
								</tr>
								{#each txs as tx}
									<tr>
										<td><a href={'#'} on:click={() => linkTx(tx.tx)}>{abbr(tx.tx)}</a></td>
										<td><a href={'#'} on:click={() => linkBlock(tx.block)}>{tx.block}</a></td>
										<td>
											{#if isCurrent(tx.from)}
												{abbr(tx.from)}
											{:else}
												<a href={'#'} on:click={() => linkAddress(tx.from)}>{abbr(tx.from)}</a>
											{/if}
										</td>
										<td>
											{#if isCurrent(tx.to)}
												{abbr(tx.to)}
											{:else}
												<a href={'#'} on:click={() => linkAddress(tx.to)}>{abbr(tx.to)}</a>
											{/if}
										</td>
										<td>{formatEth(ethers.BigNumber.from(tx.amount))}</td>
									</tr>
								{/each}
							</tbody>
						</Table>
					{/if}
				</CardBody>
			</Card>
		{/if}
	</Col>
</Row>

<style>
	.titleWidth {
		width: 10em;
	}

	.contractCode {
		font-family: monospace;
		border: 1px solid #cccccc;
		background-color: #eeeeee;
		word-wrap: break-word;
	}

	.contractCode:hover,
	.contractCode:active,
	.contractCode:focus {
		outline: 0px !important;
		-webkit-appearance: none;
		box-shadow: none !important;
	}
</style>

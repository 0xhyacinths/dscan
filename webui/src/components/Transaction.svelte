<script lang="ts">
	import {
		Col,
		Row,
		Alert,
		Spinner,
		Card,
		CardBody,
		CardHeader,
		Table,
		TabContent,
		TabPane,
		Badge
	} from 'sveltestrap';
	import { ethers } from 'ethers';
	import { onMount, afterUpdate, createEventDispatcher } from 'svelte';
	import { State } from '../lib/state';
	import type { SearchResult } from '../lib/search';
	import { ResultType } from '../lib/search';

	const dispatch = createEventDispatcher();

	function fixBN(num: ethers.BigNumber | null | undefined): ethers.BigNumber {
		return num!;
	}

	function formatGwei(num: ethers.BigNumber | null | undefined): string {
		return ethers.utils.formatUnits(num!, 'gwei');
	}

	function formatEth(num: ethers.BigNumber | null | undefined): string {
		return ethers.utils.formatUnits(num!, 'ether');
	}

	export let query: string;
	let lastQuery: string;

	let receipt: ethers.providers.TransactionReceipt;
	let txn: ethers.providers.TransactionResponse;
	let block: ethers.providers.Block;
	let state: State = State.Loading;
	let err = '';

	let senderENS: string | null;
	let receiverENS: string | null;
	let recipientContract: boolean;
	let height: number;

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
		if (window.ethereum === undefined) {
			state = State.Error;
			err = 'No Ethereum provider found.';
			return;
		}
		const provider = new ethers.providers.Web3Provider(window.ethereum, 'any');
		try {
			receipt = await provider.getTransactionReceipt(query);
			txn = await provider.getTransaction(query);
			console.log(txn);
			block = await provider.getBlock(txn.blockHash!);
			senderENS = await provider.lookupAddress(receipt.from);
			height = await provider.getBlockNumber();
			console.log(receipt);
			if (receipt.to) {
				receiverENS = await provider.lookupAddress(receipt.to);
				recipientContract = (await provider.getCode(receipt.to)).length > 2;
			}
		} catch (e) {
			err = (e as any).toString();
			state = State.Error;
			return;
		}
		state = State.Loaded;
	}

	function convertData(data: string): string {
		if (data.length <= 10) {
			return data;
		}
		const methodValue = data.slice(0, 10);
		const restValues = data.slice(10).match(/.{1,64}/g);
		let fmtData = `Method ID: ${methodValue}`;
		if (restValues) {
			const padding = restValues.length.toString().length + 1;
			return restValues
				.reduce((p, c, i) => {
					let paddingCount = padding - i.toString().length;
					return p.concat('\n').concat(`[${i}]:${' '.repeat(paddingCount)}${c}`);
				}, fmtData)
				.trim();
		}
		return fmtData.trim();
	}

	function convertLog(log: string): string {
		if (log.length <= 2) {
			return log;
		}
		const logValues = log.slice(2).match(/.{1,64}/g);
		if (logValues) {
			const padding = logValues.length.toString().length + 1;
			return logValues
				.reduce((p, c, i) => {
					let paddingCount = padding - i.toString().length;
					return p.concat('\n').concat(`[${i}]:${' '.repeat(paddingCount)}${c}`);
				}, '')
				.trim();
		}
		return log.trim();
	}

	function linkAddress(address: string) {
		dispatch('search', {
			query: address,
			type: ResultType.Address
		} as SearchResult);
	}

	function linkBlock(block: number) {
		dispatch('search', {
			query: block.toString(),
			type: ResultType.Block
		} as SearchResult);
	}
</script>

<Row>
	<Col
		>{#if state == State.Error}
			<Alert color="danger">
				<h4 class="alert-heading text-capitalize">Error loading transaction</h4>
				{err}
			</Alert>
		{:else if state == State.Loading}
			<div class="text-center">
				<Spinner color="dark" />
			</div>
		{:else if state == State.Loaded}
			<Card class="mb-3">
				<CardHeader>
					Transaction {txn.hash}
				</CardHeader>
				<CardBody>
					<TabContent>
						<TabPane tabId="overview" tab="Overview" active>
							<Table responsive borderless class="table-nomargin">
								<tbody>
									<tr>
										<th class="titleWidth" scope="row">Transaction Hash</th>
										<td>{receipt.transactionHash}</td>
									</tr>
									<tr>
										<th class="titleWidth" scope="row">Status</th>
										<td>
											{#if receipt.status == 0}
												<Badge color="danger">Failed</Badge>
											{:else if receipt.status == 1}
												<Badge color="success">Success</Badge>
											{:else}
												<Badge color="secondary">Unknown</Badge>
											{/if}
										</td>
									</tr>
									<tr>
										<th class="titleWidth" scope="row">Block</th>
										<td
											><a href={'#'} on:click={() => linkBlock(receipt.blockNumber)}
												>{receipt.blockNumber}</a
											>
											<Badge color="secondary">{height - receipt.blockNumber + 1} blocks ago</Badge
											></td
										>
									</tr>
									<tr>
										<th class="titleWidth" scope="row">Timestamp</th>
										<td>{new Date(block.timestamp * 1000).toISOString()}</td>
									</tr>
								</tbody>
							</Table>
							<hr class="hrClass" />
							<Table responsive borderless>
								<tbody>
									<tr>
										<th class="titleWidth" scope="row">From</th>
										<td
											><a href={'#'} on:click={() => linkAddress(receipt.from)}>{receipt.from}</a>
											{#if senderENS}
												({senderENS})
											{/if}
										</td>
									</tr>
									<tr>
										<th class="titleWidth" scope="row">To</th>
										{#if receipt.to}
											<td>
												<a href={'#'} on:click={() => linkAddress(receipt.to)}>{receipt.to}</a>
												{#if receiverENS}
													({receiverENS})
												{/if}
												{#if recipientContract}
													(Contract)
												{/if}</td
											>
										{:else if receipt.contractAddress}
											<td>
												<a href={'#'} on:click={() => linkAddress(receipt.contractAddress)}
													>{receipt.contractAddress}</a
												>
												(Contract Created)</td
											>
										{/if}
									</tr>
								</tbody>
							</Table>
							<hr class="hrClass" />
							<Table responsive borderless>
								<tbody>
									<tr>
										<th class="titleWidth" scope="row">Value</th>
										<td>{formatEth(txn.value)} Ether</td>
									</tr>
									<tr>
										<th class="titleWidth" scope="row">Transaction Fee</th>
										<td>{formatEth(receipt.effectiveGasPrice.mul(receipt.gasUsed))} Ether</td>
									</tr>
									{#if receipt.to}
										<tr>
											<th class="titleWidth" scope="row">Data</th>
											<td><pre class="noMarginBottom">{convertData(txn.data)}</pre></td>
										</tr>
									{:else if receipt.contractAddress}
										<tr>
											<th class="titleWidth" scope="row">Creation Code</th>
											<td
												><textarea
													class="form-control contractCode"
													rows="10"
													name="text"
													readonly
													bind:value={txn.data}
												/></td
											>
										</tr>
									{/if}
								</tbody>
							</Table>
							<hr class="hrClass" />
							<Table responsive borderless>
								<tbody>
									<tr>
										<th class="titleWidth" scope="row">Gas Price</th>
										<td>{formatGwei(receipt.effectiveGasPrice)} Gwei</td>
									</tr>
									<tr>
										<th class="titleWidth" scope="row">Gas Utilization</th>
										<td>{receipt.gasUsed} / {txn.gasLimit}</td>
									</tr>
									{#if txn.type === 2}
										<tr>
											<th class="titleWidth" scope="row">Max Fee</th>
											<td>{formatGwei(txn.maxFeePerGas)} Gwei</td>
										</tr>
										<tr>
											<th class="titleWidth" scope="row">Base Fee</th>
											<td>{formatGwei(block.baseFeePerGas)} Gwei</td>
										</tr>
										<tr>
											<th class="titleWidth" scope="row">Max Priority Fee</th>
											<td>{formatGwei(txn.maxPriorityFeePerGas)} Gwei</td>
										</tr>
										<tr>
											<th class="titleWidth" scope="row">EIP-1559 Burned</th>
											<td>{formatEth(fixBN(block.baseFeePerGas).mul(receipt.gasUsed))} Ether</td>
										</tr>
									{/if}
								</tbody>
							</Table>
						</TabPane>
						<TabPane tabId="logs" tab="Logs ({receipt.logs.length})">
							{#each receipt.logs as log, idx}
								<Table responsive borderless>
									<tbody>
										<tr>
											<th class="titleWidth" scope="row">Index</th>
											<td>{log.logIndex}</td>
										</tr>
										<tr>
											<th class="titleWidth" scope="row">Address</th>

											<td
												><a href={'#'} on:click={() => linkAddress(log.address)}>{log.address}</a
												></td
											>
										</tr>
										<tr>
											<th class="titleWidth" scope="row">Data</th>
											<td><pre class="noMarginBottom">{convertLog(log.data)}</pre></td>
										</tr>
										{#each log.topics as topic, num}
											<tr>
												<th class="titleWidth" scope="row">Topic [{num}]</th>
												<td>{topic}</td>
											</tr>
										{/each}
									</tbody>
								</Table>
								{#if idx < receipt.logs.length - 1}
									<hr class="hrClass" />
								{/if}
							{/each}
						</TabPane>
					</TabContent>
				</CardBody>
			</Card>
		{/if}
	</Col>
</Row>

<style>
	.titleWidth {
		width: 10em;
	}

	.hrClass {
		margin-top: 0px;
	}
	.noMarginBottom {
		margin-bottom: 0px;
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

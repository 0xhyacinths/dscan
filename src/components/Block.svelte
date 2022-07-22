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
		Collapse
	} from 'sveltestrap';
	import { ethers } from 'ethers';
	import { onMount, afterUpdate, createEventDispatcher } from 'svelte';
	import { State } from '../lib/state';

	import type { SearchResult } from '../lib/search';
	import { ResultType } from '../lib/search';

	const dispatch = createEventDispatcher();

	export let query: string;
	let lastQuery: string;

	let block: ethers.providers.Block;
	let state: State = State.Loading;
	let err = '';
	let minerENS: string | null;
	let expanded = false;

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
			if (query.length == 66) {
				block = await provider.getBlock(query);
			} else {
				block = await provider.getBlock(parseInt(query));
			}
			minerENS = await provider.lookupAddress(block.miner);
		} catch (e) {
			err = (e as any).toString();
			state = State.Error;
			return;
		}
		state = State.Loaded;
	}

	function formatGwei(num: ethers.BigNumber | null | undefined): string {
		return ethers.utils.formatUnits(num!, 'gwei');
	}

	function linkAddress(address: string) {
		dispatch('search', {
			query: address,
			type: ResultType.Address
		} as SearchResult);
	}

	function linkTx(tx: string) {
		dispatch('search', {
			query: tx,
			type: ResultType.Transaction
		} as SearchResult);
	}
</script>

<Row>
	<Col
		>{#if state == State.Error}
			<Alert color="danger">
				<h4 class="alert-heading text-capitalize">Error loading block</h4>
				{err}
			</Alert>
		{:else if state == State.Loading}
			<div class="text-center">
				<Spinner color="dark" />
			</div>
		{:else if state == State.Loaded}
			<Card class="mb-3">
				<CardHeader>
					Block {query}
				</CardHeader>
				<CardBody>
					<Table responsive borderless class="table-nomargin">
						<tbody>
							<tr>
								<th class="titleWidth" scope="row">Number</th>
								<td>{block.number}</td>
							</tr>
							<tr>
								<th class="titleWidth" scope="row">Hash</th>
								<td>{block.hash}</td>
							</tr>
							<tr>
								<th class="titleWidth" scope="row">Timestamp</th>
								<td>{new Date(block.timestamp * 1000).toISOString()}</td>
							</tr>
							<tr>
								<th class="titleWidth" scope="row">Transactions</th>
								<td>
									{block.transactions.length}
									{#if block.transactions.length > 0}
										<a href={'#'} on:click={() => (expanded = !expanded)}>
											{#if expanded}
												Hide
											{:else}
												Show
											{/if}
										</a>
										<Collapse isOpen={expanded}>
											<div class="big">
												{#each block.transactions as transaction, idx}
													{idx + 1}:
													<a href={'#'} on:click={() => linkTx(transaction)}>{transaction}</a><br />
												{/each}
											</div>
										</Collapse>
									{/if}
								</td>
							</tr>
							<tr>
								<th class="titleWidth" scope="row">Miner</th>
								<td
									><a href={'#'} on:click={() => linkAddress(block.miner)}>{block.miner}</a>
									{#if minerENS}({minerENS}){/if}</td
								>
							</tr>
							<tr>
								<th class="titleWidth" scope="row">Gas Utilization</th>
								<td>{block.gasUsed.toString()} / {block.gasLimit.toString()}</td>
							</tr>
							{#if block.baseFeePerGas}
								<tr>
									<th class="titleWidth" scope="row">Base Fee</th>
									<td>{formatGwei(block.baseFeePerGas)} Gwei</td>
								</tr>
							{/if}
							<tr>
								<th class="titleWidth" scope="row">Difficulty</th>
								<td>{block._difficulty.toString()}</td>
							</tr>
							<tr>
								<th class="titleWidth" scope="row">Extra data</th>
								<td
									><textarea
										class="form-control extraData"
										rows="2"
										name="text"
										readonly
										bind:value={block.extraData}
									/></td
								>
							</tr>
						</tbody>
					</Table>
				</CardBody>
			</Card>
		{/if}
	</Col>
</Row>

<style>
	.titleWidth {
		width: 10em;
	}

	.extraData {
		font-family: monospace;
		border: 1px solid #cccccc;
		background-color: #eeeeee;
		word-wrap: break-word;
	}

	.extraData:hover,
	.extraData:active,
	.extraData:focus {
		outline: 0px !important;
		-webkit-appearance: none;
		box-shadow: none !important;
	}

	.big {
		max-height: 200px;
		overflow: scroll;
	}
</style>

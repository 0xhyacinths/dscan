<script lang="ts">
	import { page } from '$app/stores';
	import {
		Col,
		Container,
		Row,
		Alert,
		Spinner,
		Card,
		CardBody,
		CardHeader,
		Table
	} from 'sveltestrap';
	import { ethers } from 'ethers';
	import { onMount } from 'svelte';

	enum State {
		Loading,
		Error,
		Loaded
	}

	function fixBN(num: ethers.BigNumber | null | undefined): ethers.BigNumber {
		return num!;
	}

	function formatGwei(num: ethers.BigNumber | null | undefined): string {
		return ethers.utils.formatUnits(num!, 'gwei');
	}

	function formatEth(num: ethers.BigNumber | null | undefined): string {
		return ethers.utils.formatUnits(num!, 'ether');
	}

	const id = $page.url.searchParams.get('id');
	let receipt: ethers.providers.TransactionReceipt;
	let txn: ethers.providers.TransactionResponse;
	let block: ethers.providers.Block;
	let state: State = State.Loading;
	let err = '';

	let minerENS: string | null;

	onMount(async () => {
		if (window.ethereum === undefined) {
			state = State.Error;
			err = 'No Ethereum provider found.';
			return;
		}
		const provider = new ethers.providers.Web3Provider(window.ethereum, 'any');
		try {
			if (id!.length == 66) {
				block = await provider.getBlock(id!);
			} else {
				block = await provider.getBlock(parseInt(id!));
			}
			minerENS = await provider.lookupAddress(block.miner);
		} catch (e) {
			err = (e as any).toString();
			state = State.Error;
			return;
		}
		state = State.Loaded;
	});

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
</script>

<svelte:head>
	<title>Block {id}</title>
</svelte:head>

<Container>
	<div class="padded" />
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
						Block {id}
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
									<td>{block.transactions.length}</td>
								</tr>
								<tr>
									<th class="titleWidth" scope="row">Miner</th>
									<td
										><a href="/address?id={block.miner}">{block.miner}</a>
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
</Container>

<style>
	.padded {
		margin-top: 2rem;
	}

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
</style>

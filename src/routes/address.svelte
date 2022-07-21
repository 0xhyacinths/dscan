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

	function formatEth(num: ethers.BigNumber | null | undefined): string {
		return ethers.utils.formatUnits(num!, 'ether');
	}

	const id = $page.url.searchParams.get('id');
	let state: State = State.Loading;
	let err = '';

	let balance: ethers.BigNumber;
	let contractCode: string;
	let totalTx: number;
	let ensAddress: string | null;

	onMount(async () => {
		if (window.ethereum === undefined) {
			state = State.Error;
			err = 'No Ethereum provider found.';
			return;
		}
		const provider = new ethers.providers.Web3Provider(window.ethereum, 'any');
		try {
			balance = await provider.getBalance(id!);
			contractCode = await provider.getCode(id!);
			totalTx = await provider.getTransactionCount(id!);
			ensAddress = await provider.lookupAddress(id!);
		} catch (e) {
			err = (e as any).toString();
			state = State.Error;
			return;
		}
		state = State.Loaded;
	});
</script>

<svelte:head>
	<title>Address {id}</title>
</svelte:head>

<Container>
	<div class="padded" />
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
						Address {id}
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

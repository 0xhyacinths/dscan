<script lang="ts">
	import { page } from '$app/stores';
	import { Col, Container, Row } from 'sveltestrap';
	import Nav from '../components/Nav.svelte';
	import { onMount } from 'svelte';
	import Address from '../components/Address.svelte';
	import Block from '../components/Block.svelte';
	import Transaction from '../components/Transaction.svelte';
	import type { SearchResult } from '../lib/search';
	import { ResultType } from '../lib/search';
	import { goto } from '$app/navigation';
	import { DescanClient } from '../lib/api';
	import { ethers } from 'ethers';

	const version: string = __APP_VERSION__;
  const basedir: string = __BASE_PATH__;
	$: current = getCurrent($page.url.searchParams);

	function getCurrent(params: URLSearchParams): SearchResult | null {
		const typeStr = params.get('type');
		const queryStr = params.get('query');
		if (typeStr && queryStr) {
			const type: ResultType = parseInt(typeStr) as ResultType;
			return {
				query: queryStr,
				type: type
			} as SearchResult;
		}
		return null;
	}

	let client: DescanClient = new DescanClient('http://127.0.0.1:9090', 0);

	onMount(async () => {
		const endpoint = window.localStorage.getItem('endpoint') ?? 'http://127.0.0.1:9090';
		client.setEndpoint(endpoint);
		if (window.ethereum !== undefined) {
			const provider = new ethers.providers.Web3Provider(window.ethereum, 'any');
			const network = await provider.getNetwork();
			client.setNetwork(network.chainId);
		}
		current = getCurrent($page.url.searchParams);
		console.log('mount');
	});

	async function handleSearch(res: CustomEvent<SearchResult>) {
		await updateSearch(res);
		console.log('gone to');
		current = getCurrent($page.url.searchParams);
	}

	async function updateSearch(res: CustomEvent<SearchResult>) {
		let qs = new URLSearchParams($page.url.searchParams.toString());
		qs.set('type', res.detail.type.toString());
		qs.set('query', res.detail.query);
		console.log(qs.toString());
		await goto(`${basedir}/?${qs.toString()}`);
	}
</script>

<svelte:head>
	<title>Decentralized Explorer</title>
</svelte:head>

<Nav on:search={handleSearch} base={basedir} {client} />
<Container>
	<div class="padded" />
	<Row>
		{#if current}
			{#if current.type == ResultType.Address}
				<Address query={current} on:search={handleSearch} on:updateSearch={updateSearch} {client} />
			{:else if current.type == ResultType.Block}
				<Block query={current.query} on:search={handleSearch} />
			{:else if current.type == ResultType.Transaction}
				<Transaction query={current.query} on:search={handleSearch} />
			{/if}
		{:else}
			<Row>
				<Col
					><h1>DScan: A Decentralized Explorer</h1>
					<p>DScan provides a minimally-centralized Ethereum blockchain explorer.</p></Col
				>
			</Row>
		{/if}
	</Row>
</Container>

<div class="container">
	<footer class="py-3 my-4">
		<p class="text-center text-muted">
			<a href="https://github.com/0xhyacinths/dscan" class="text-muted">DScan</a>
			{version}
		</p>
	</footer>
</div>

<style>
	.padded {
		margin-top: 2rem;
	}
</style>


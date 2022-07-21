<script lang="ts">
	import { page } from '$app/stores';
	import { Col, Container, Row } from 'sveltestrap';
	import Nav from '../components/Nav.svelte';
  import {onMount} from 'svelte';

	import Address from '../components/Address.svelte';
	import Block from '../components/Block.svelte';
	import Transaction from '../components/Transaction.svelte';

	import type { SearchResult } from '../lib/search';
	import { ResultType } from '../lib/search';
	import { goto } from '$app/navigation';

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

  onMount(() => {
    console.log("mount");
  });

	function handleSearch(res: CustomEvent<SearchResult>) {
		let qs = new URLSearchParams($page.url.searchParams.toString());
		qs.set('type', res.detail.type.toString());
		qs.set('query', res.detail.query);
		console.log(qs.toString());
		goto(`/?${qs.toString()}`);
	}
</script>

<svelte:head>
	<title>Decentralized Explorer</title>
</svelte:head>

<Nav on:search={handleSearch} />
<Container>
	<div class="padded" />
	<Row>
		{#if current}
			{#if current.type == ResultType.Address}
				<Address query={current.query} />
			{:else if current.type == ResultType.Block}
				<Block query={current.query} on:search={handleSearch} />
			{:else if current.type == ResultType.Transaction}
				<Transaction query={current.query} on:search={handleSearch} />
			{/if}
		{:else}
			<Row>
				<Col><h1>Decentralized Explorer</h1></Col>
			</Row>
		{/if}
	</Row>
</Container>

<style>
	.padded {
		margin-top: 2rem;
	}
</style>

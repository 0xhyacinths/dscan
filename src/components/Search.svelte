<script lang="ts">
	import type { SearchResult } from '../lib/search';
	import { ResultType } from '../lib/search';
	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();

	let query: string;
	function handleSearch() {
		switch (true) {
			case /0x[0-9a-fA-F]{64}/.test(query):
				dispatch('search', {
					query: query,
					type: ResultType.Transaction
				} as SearchResult);
				break;
			case /0x[0-9a-fA-F]{40}/.test(query):
				dispatch('search', {
					query: query,
					type: ResultType.Address
				} as SearchResult);
				break;
      case /.*\.eth/.test(query):
				dispatch('search', {
					query: query,
					type: ResultType.Address
				} as SearchResult);
				break;
			case /[0-9]+/.test(query):
				dispatch('search', {
					query: query,
					type: ResultType.Block
				} as SearchResult);
				break;
			default:
				break;
		}
	}
</script>

<form on:submit|preventDefault={handleSearch}>
	<div class="input-group input-group-lg big-width">
		<input
			type="text"
			class="form-control"
			placeholder="Address, block number, or transaction hash"
			aria-label="Address, block number, or transaction hash"
			aria-describedby="basic-addon2"
			bind:value={query}
		/>
		<div class="input-group-append">
			<button class="btn btn-lg btn-primary" type="button" on:click={handleSearch}>Search</button>
		</div>
	</div>
</form>

<style>
	.big-width {
		width: 50vw;
	}
</style>

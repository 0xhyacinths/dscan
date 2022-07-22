<script lang="ts">
	import { Navbar, NavbarBrand, Button, Popover } from 'sveltestrap';
	import type { DescanClient } from '../lib/api';
	import Search from './Search.svelte';

	export let client: DescanClient;
	let alive: boolean = false;
	let message: string = '';

	client.alive.subscribe((val) => {
		alive = val;
	});
	client.message.subscribe((val) => {
		message = val;
	});
</script>

<Navbar color="light" light expand="md">
	<NavbarBrand href="/">DScan</NavbarBrand>
	<Search on:search />
	<Button id="btn-bottom">
		{#if alive}
			Connected
		{:else}
			Disconnected
		{/if}
	</Button>
	<Popover target="btn-bottom" placement="bottom" title="Connection Status">
		{#if alive}
			Connected to {client.getEndpoint()}
			<hr />
			{message}
		{:else}
			Not connected
		{/if}
	</Popover>
</Navbar>
